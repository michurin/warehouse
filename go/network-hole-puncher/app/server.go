package app

import (
	"errors"
	"fmt"
	"net"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/udp"
)

var noDataErr = errors.New("Empty data")

func newServerHandler() udp.Handler {
	addresses := [][]byte{nil, nil}
	return func(
		conn *net.UDPConn,
		addr *net.UDPAddr,
		data []byte,
	) (bool, error) {
		if len(data) < 1 {
			return false, noDataErr
		}
		idx := int(data[0]) & 1
		addresses[idx] = append([]byte("PEER@"), append(append(data, '@'), []byte(addr.String())...)...)
		payload := addresses[idx^1]
		err := udp.Send(conn, addr, payload)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		return false, nil
	}
}

func Server(address string) error {
	conn, err := udp.Connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = udp.Serve(conn, newServerHandler())
	if err != nil {
		return err
	}
	return nil
}
