package app

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/udp"
)

const ( // TODO move to sep file
	labelPeerInfo   = 'I'
	labelPing       = 'x'
	labelPong       = 'y'
	labelClose      = 'z'
	labelsSeporator = '|'
)

var noDataErr = errors.New("Empty data")

func newServerHandler() udp.Handler {
	addresses := [][]byte{nil, nil}
	return func(
		conn *net.UDPConn,
		addr *net.UDPAddr,
		data []byte,
	) {
		if len(data) < 1 {
			return
		}
		idx := int(data[0]) & 1
		addresses[idx] = bytes.Join([][]byte{
			{labelPeerInfo},
			data,
			[]byte(addr.String()),
		}, []byte{labelsSeporator})
		payload := addresses[idx^1]
		err := udp.Send(conn, addr, payload)
		if err != nil {
			fmt.Println(err)
			return
		}
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
