package proxyproto

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"strings"
	"time"
)

// addr -- implements net.Addr
type addr struct {
	net  string
	addr string
}

func (a *addr) Network() string { return a.net }
func (a *addr) String() string  { return a.addr }

// proxyConn -- implements net.Conn
type proxyConn struct {
	net.Conn

	buf *bufio.Reader

	localAddr  net.Addr
	remoteAddr net.Addr

	err error
}

func (c *proxyConn) Close() error {
	return c.Conn.Close()
}

func (c *proxyConn) LocalAddr() net.Addr { return c.localAddr }

func (c *proxyConn) Read(buff []byte) (n int, err error) {
	if c.err != nil {
		return 0, c.err
	}
	return c.buf.Read(buff)
}

func (c *proxyConn) RemoteAddr() net.Addr { return c.remoteAddr }

// ProxyListener -- implements net.Listener but requires each connection to
// start with a valid PROXY line
// see https://blog.digitalocean.com/load-balancers-now-support-proxy-protocol/
type ProxyListener struct {
	raw net.Listener
}

// Accept -- waits for incoming connections
func (l *ProxyListener) Accept() (net.Conn, error) {
	var raw, err = l.raw.Accept()
	if err != nil {
		return raw, err
	}

	if tcpConn, ok := raw.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(3 * time.Minute)
	}

	raw.SetDeadline(time.Now().Add(10 * time.Second))
	var buf = bufio.NewReader(raw)

	var localAddr, remoteAddr net.Addr
	if localAddr, remoteAddr, err = l.parseProxyLine(buf); err != nil {
		return &proxyConn{
			Conn:       raw,
			buf:        buf,
			localAddr:  raw.LocalAddr(),
			remoteAddr: raw.RemoteAddr(),
			err:        err,
		}, nil
	}

	raw.SetDeadline(time.Time{})
	return &proxyConn{
		Conn:       raw,
		buf:        buf,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}, nil
}

// Addr -- returns the address this Listener listens at
func (l *ProxyListener) Addr() net.Addr {
	return l.raw.Addr()
}

// Close -- Stops listening
func (l *ProxyListener) Close() error {
	return l.raw.Close()
}

func (l *ProxyListener) parseProxyLine(buf *bufio.Reader) (localAddr, remoteAddr net.Addr, err error) {
	// TODO set a strict timeout here - otherwise we could end up slowing down other requests

	// format (v1) "PROXY_STRING INET_PROTOCOL CLIENT_IP PROXY_IP CLIENT_PORT PROXY_PORT\r\n"
	var proxyWord []byte
	if proxyWord, err = buf.Peek(6); err != nil {
		return nil, nil, err
	}
	if !bytes.Equal(proxyWord, []byte("PROXY ")) {
		return nil, nil, errors.New("not a PROXY request")
	}

	var line string
	if line, err = buf.ReadString('\n'); err != nil {
		return nil, nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	//logrus.Info("proxy line: ", line)
	var parts = strings.Split(line, " ")
	if len(parts) != 6 {
		return nil, nil, errors.New("malformed PROXY line")
	}

	var clientIP, proxyIP, clientPort, proxyPort = parts[2], parts[3], parts[4], parts[5]
	localAddr = &addr{"tcp", net.JoinHostPort(proxyIP, proxyPort)}
	remoteAddr = &addr{"tcp", net.JoinHostPort(clientIP, clientPort)}

	return localAddr, remoteAddr, nil
}

// NewProxyListener -- creates and returns a ProxyListener (wrapping the given net.Listener)
func NewProxyListener(raw net.Listener) net.Listener {
	var rc = ProxyListener{
		raw: raw,
	}
	return &rc
}
