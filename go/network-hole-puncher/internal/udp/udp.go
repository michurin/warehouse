package udp

import (
	"fmt"
	"net"
)

type Handler func(*net.UDPConn, *net.UDPAddr, []byte) (bool, error)

func Connect(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Serve(
	conn *net.UDPConn,
	handler Handler,
) error {
	for {
		fmt.Println("reading...")
		data := make([]byte, 1024) // we create new slice every time to prevent sharing memory betwinn server and handler
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			return err
		}
		fmt.Printf("call handler: data=%q, from %s\n", string(data[:n]), addr)
		last, err := handler(conn, addr, data[:n])
		if err != nil {
			return nil
		}
		if last {
			break
		}
	}
	return nil
}

func Send(conn *net.UDPConn, addr *net.UDPAddr, data []byte) error {
	_, err := conn.WriteToUDP(data, addr)
	return err
}
