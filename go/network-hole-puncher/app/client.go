package app

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/udp"
)

var slotTooLongError = errors.New("Slot is too long")

type peer struct {
	mx      *sync.Mutex
	address net.UDPAddr
	ttl     int
}

func newPeer() *peer {
	return &peer{
		mx: new(sync.Mutex),
	}
}

func (p *peer) Get() (net.UDPAddr, bool) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.ttl <= 0 {
		return net.UDPAddr{}, false
	}
	p.ttl--
	return p.address, true
}

func (p *peer) Set(addr net.UDPAddr) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.address = addr
	p.ttl = 10
}

func newClientHandler(p *peer) udp.Handler {
	return func(
		conn *net.UDPConn,
		addr *net.UDPAddr,
		data []byte,
	) (bool, error) {
		fmt.Printf("THE RESULT: REMOTE GLOBAL ADDRESS: %q\n", data)
		// TODO p.Set() if it is data from server
		// TODO reply if it is data from peer
		// TODO finish if it is PONG from peer (Nth pong?)
		return true, nil // TODO exit only if we obtain address
	}
}

func Client(slot, address, remoteAddress string) error {
	if len(slot) > 16 {
		return slotTooLongError
	}
	message := []byte(slot)

	p := newPeer()

	conn, err := udp.Connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	addr, err := net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		return err
	}

	go func() { // TODO repeat sending
		// TODO p.Get()
		// TODO - send to Server
		// TODO - or send to Peer if Peer available
		fmt.Println("Sending to server...")
		err = udp.Send(conn, addr, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	err = udp.Serve(conn, newClientHandler(p))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server is stopped")
	return nil
}

