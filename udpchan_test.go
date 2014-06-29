package udpchan

import (
	"bytes"
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatal(err)
	}

	udpC, err := Connect(":9999")
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("foo")
	udpC <- message

	read := make([]byte, mtu)
	n, err := conn.Read(read)
	if err != nil {
		t.Fatal(err)
	}
	read = read[:n]

	if bytes.Compare(read, message) != 0 {
		t.Errorf("Expected result %v, got %v", message, read)
	}

	conn.Close()
}

func TestConnectAndListen(t *testing.T) {
	done := make(chan bool)

	inbound, err := Listen(":9999", done)
	if err != nil {
		t.Fatal(err)
	}

	outbound, err := Connect(":9999")
	if err != nil {
		t.Fatal(err)
	}

	message := []byte("foo")
	outbound <- message
	read := <-inbound
	if bytes.Compare(read, message) != 0 {
		t.Errorf("Expected result %v, got %v", message, read)
	}

	// Probably blocked on a channel read,
	// so just in case.
	outbound <- []byte("")
	done <- true
}

// Just make sure we're closing the listener
func TestConnectAndListen2(t *testing.T) {
	TestConnectAndListen(t)
}
