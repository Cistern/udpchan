// Package udpchan provides a thin channel wrapper
// around UDP connections.
package udpchan

import (
	"net"
)

// Usually the MTU is 1500 bytes and
// a little over 10k for jumbo frames.
// 10k should be plenty.
const mtu = 10000

// Connect dials to address and returns a send-only
// channel and an error.
func Connect(address string) (chan<- []byte, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	outbound := make(chan []byte, 10)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Caught a panic, so we're going to close
				// the outbound channel. Receivers are supposed
				// to check for closed channels, so we're good.
				close(outbound)
			}
		}()
		for message := range outbound {
			conn.Write(message)
		}
	}()

	return outbound, nil
}

// Listen starts a UDP listener at address and
// returns a read-only channel and an error.
// If a value is sent to close, the UDP socket will
// be closed.
func Listen(address string, close chan bool) (<-chan []byte, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":9999")
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}

	inbound := make(chan []byte, 10)
	go func() {
		for {
			select {
			case <-close:
				conn.Close()
				return
			default:
			}

			b := make([]byte, mtu)
			n, err := conn.Read(b)
			if err != nil {
				continue
			}

			b = b[:n]
			inbound <- b
		}
	}()

	return inbound, nil
}
