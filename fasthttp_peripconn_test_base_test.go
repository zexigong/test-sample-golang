package fasthttp

import (
	"crypto/tls"
	"net"
	"testing"
)

type mockConn struct {
	closed bool
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP: net.IPv4(192, 168, 1, 1),
	}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestPerIPConnCounter_RegisterUnregister(t *testing.T) {
	counter := &perIPConnCounter{}
	ip := uint32(3232235777) // 192.168.1.1

	n := counter.Register(ip)
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}

	n = counter.Register(ip)
	if n != 2 {
		t.Fatalf("expected 2, got %d", n)
	}

	counter.Unregister(ip)
	if counter.m[ip] != 1 {
		t.Fatalf("expected 1, got %d", counter.m[ip])
	}

	counter.Unregister(ip)
	if counter.m[ip] != 0 {
		t.Fatalf("expected 0, got %d", counter.m[ip])
	}
}

func TestAcquirePerIPConn(t *testing.T) {
	counter := &perIPConnCounter{}
	ip := uint32(3232235777) // 192.168.1.1
	conn := &mockConn{}

	perIPConn := acquirePerIPConn(conn, ip, counter)

	if perIPConn == nil {
		t.Fatal("expected non-nil perIPConn")
	}

	if _, ok := perIPConn.(*perIPConn); !ok {
		t.Fatalf("expected *perIPConn, got %T", perIPConn)
	}

	if err := perIPConn.Close(); err != nil {
		t.Fatalf("unexpected error on close: %v", err)
	}

	if !conn.closed {
		t.Fatalf("expected conn to be closed")
	}
}

func TestAcquirePerIPTLSConn(t *testing.T) {
	counter := &perIPConnCounter{}
	ip := uint32(3232235777) // 192.168.1.1
	tlsConn := tls.Client(&mockConn{}, &tls.Config{})

	perIPConn := acquirePerIPConn(tlsConn, ip, counter)

	if perIPConn == nil {
		t.Fatal("expected non-nil perIPTLSConn")
	}

	if _, ok := perIPConn.(*perIPTLSConn); !ok {
		t.Fatalf("expected *perIPTLSConn, got %T", perIPConn)
	}

	if err := perIPConn.Close(); err != nil {
		t.Fatalf("unexpected error on close: %v", err)
	}
}

func TestGetUint32IP(t *testing.T) {
	conn := &mockConn{}
	expectedIP := uint32(3232235777) // 192.168.1.1

	ip := getUint32IP(conn)
	if ip != expectedIP {
		t.Fatalf("expected %d, got %d", expectedIP, ip)
	}
}

func TestIPConversion(t *testing.T) {
	originalIP := net.IPv4(192, 168, 1, 1)
	uint32IP := ip2uint32(originalIP)
	convertedIP := uint322ip(uint32IP)

	if !originalIP.Equal(convertedIP) {
		t.Fatalf("expected %v, got %v", originalIP, convertedIP)
	}
}