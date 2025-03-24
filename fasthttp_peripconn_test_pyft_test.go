package fasthttp

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestPerIPConnCounter(t *testing.T) {
	t.Parallel()

	var cc perIPConnCounter

	ip := ip2uint32(net.ParseIP("1.2.3.4"))

	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}
	n = cc.Register(ip)
	if n != 2 {
		t.Fatalf("unexpected n: %d. Expecting 2", n)
	}

	cc.Unregister(ip)
	n = cc.Register(ip)
	if n != 2 {
		t.Fatalf("unexpected n: %d. Expecting 2", n)
	}

	cc.Unregister(ip)
	cc.Unregister(ip)
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}
}

func TestPerIPConn(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	ch := make(chan struct{})
	go func() {
		defer close(ch)

		c, err := ln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	<-ch
}

func TestPerIPConnTLS(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair("./testdata/ssl-cert-snakeoil.pem", "./testdata/ssl-cert-snakeoil.key")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tln := tls.NewListener(ln, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	go func() {
		c, err := tln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := tls.Dial("tcp4", ln.Addr().String(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}
}

func TestGetConnIP4(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	ipCh := make(chan net.IP, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		ipCh <- getConnIP4(c)
	}()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer c.Close()

	select {
	case ip := <-ipCh:
		ipExpected := net.ParseIP("127.0.0.1").To4()
		if !ip.Equal(ipExpected) {
			t.Fatalf("unexpected ip: %s. Expecting %s", ip, ipExpected)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
}

func TestIP2Uint32(t *testing.T) {
	t.Parallel()

	ip := net.ParseIP("10.20.30.40").To4()
	ipUint32 := ip2uint32(ip)
	if ipUint32 != 0x0a141e28 {
		t.Fatalf("unexpected ipUint32: %08x. Expecting 0x0a141e28", ipUint32)
	}
}

func TestUint322IP(t *testing.T) {
	t.Parallel()

	ip := uint322ip(0x0a141e28)
	expectedIP := net.ParseIP("10.20.30.40").To4()
	if !ip.Equal(expectedIP) {
		t.Fatalf("unexpected ip: %s. Expecting %s", ip, expectedIP)
	}
}

type errorAddr struct{}

func (errorAddr) Network() string {
	return "network"
}

func (errorAddr) String() string {
	return "string"
}

type errorConn struct{}

func (errorConn) Read(b []byte) (int, error) {
	return 0, errors.New("errorConn")
}

func (errorConn) Write(b []byte) (int, error) {
	return 0, errors.New("errorConn")
}

func (errorConn) Close() error {
	return errors.New("errorConn")
}

func (errorConn) LocalAddr() net.Addr {
	return errorAddr{}
}

func (errorConn) RemoteAddr() net.Addr {
	return errorAddr{}
}

func (errorConn) SetDeadline(t time.Time) error {
	return errors.New("errorConn")
}

func (errorConn) SetReadDeadline(t time.Time) error {
	return errors.New("errorConn")
}

func (errorConn) SetWriteDeadline(t time.Time) error {
	return errors.New("errorConn")
}

func TestGetConnIP4Error(t *testing.T) {
	t.Parallel()

	ip := getConnIP4(errorConn{})
	if ip != nil {
		t.Fatalf("ip must be nil")
	}

	ip = getConnIP4(nil)
	if ip != nil {
		t.Fatalf("ip must be nil")
	}
}

func TestGetUint32IP(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	ipCh := make(chan uint32, 1)
	go func() {
		c, err := ln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		ipCh <- getUint32IP(c)
	}()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer c.Close()

	select {
	case ipUint32 := <-ipCh:
		expectedIPUint32 := ip2uint32(net.ParseIP("127.0.0.1").To4())
		if ipUint32 != expectedIPUint32 {
			t.Fatalf("unexpected ipUint32: %08x. Expecting %08x", ipUint32, expectedIPUint32)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timeout")
	}
}

func TestGetUint32IPError(t *testing.T) {
	t.Parallel()

	ipUint32 := getUint32IP(errorConn{})
	if ipUint32 != 0 {
		t.Fatalf("unexpected ipUint32: %08x. Expecting 0", ipUint32)
	}

	ipUint32 = getUint32IP(nil)
	if ipUint32 != 0 {
		t.Fatalf("unexpected ipUint32: %08x. Expecting 0", ipUint32)
	}
}

func TestAcquirePerIPConn(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	ch := make(chan struct{})
	go func() {
		defer close(ch)

		c, err := ln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	<-ch
}

func TestAcquirePerIPConnTLS(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair("./testdata/ssl-cert-snakeoil.pem", "./testdata/ssl-cert-snakeoil.key")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tln := tls.NewListener(ln, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	go func() {
		c, err := tln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := tls.Dial("tcp4", ln.Addr().String(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}
}

type testConn struct {
	net.Conn
}

func (c *testConn) Close() error {
	return io.EOF
}

func TestPerIPConnCloseError(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	ch := make(chan struct{})
	go func() {
		defer close(ch)

		c, err := ln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	c, err := net.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	c = &testConn{c}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	<-ch
}

func TestPerIPTLSConnCloseError(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	defer ln.Close()

	cert, err := tls.LoadX509KeyPair("./testdata/ssl-cert-snakeoil.pem", "./testdata/ssl-cert-snakeoil.key")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	tln := tls.NewListener(ln, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	go func() {
		c, err := tln.Accept()
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		defer c.Close()

		b := make([]byte, 1)
		if _, err = c.Read(b); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}
		if b[0] != 42 {
			t.Errorf("unexpected b[0]: %d. Expecting 42", b[0])
			return
		}

		if err = c.Close(); err != nil {
			t.Errorf("unexpected error: %s", err)
			return
		}

		if _, err = c.Read(b); err == nil {
			t.Errorf("expecting error")
			return
		}
	}()

	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, err := tls.Dial("tcp4", ln.Addr().String(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	c = &testConn{c}

	var cc perIPConnCounter
	ip := getUint32IP(c)
	n := cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	pc := acquirePerIPConn(c, ip, &cc)
	defer pc.Close()

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}

	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}

	if _, err = pc.Write([]byte{42}); err == nil {
		t.Fatalf("expecting error")
	}

	if err = pc.Close(); err != io.EOF {
		t.Fatalf("unexpected error: %s", err)
	}
	n = cc.Register(ip)
	if n != 1 {
		t.Fatalf("unexpected n: %d. Expecting 1", n)
	}
}