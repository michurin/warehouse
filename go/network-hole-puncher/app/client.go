package app

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/udp"
)

var slotTooLongError = errors.New("Slot is too long")

type peer struct { // TODO move it to separate file
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
		fmt.Println(data)
		ff := bytes.Split(data, []byte{'@'})
		fmt.Println(ff)
		if len(ff) == 0 { // TODO we have to have this check; HOWEVER server has to return correct signature
			return false, nil
		}
		if len(ff[0]) == 0 {
			return false, nil
		}
		if ff[0][1] == 'E' { // (PEER) we receive peer information
			addr, err := net.ResolveUDPAddr("udp", string(ff[2]))
			if err != nil {
				return false, err
			}
			p.Set(*addr)
			udp.Send(conn, addr, []byte("PING@0")) // TODO few tries
			return false, nil                      // keep going
		}
		if ff[0][1] == 'I' { // (PING) from host
		}
		// TODO p.Set() if it is data from server
		// TODO reply if it is data from peer
		// TODO finish if it is PONG from peer (Nth pong?)
		return true, nil // TODO exit only if we obtain address
	}
}

func Client(slot, address, remoteAddress string) error {
	// TODO check slot letters and digits only!
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
		for {
			peer, ok := p.Get()
			fmt.Println(peer, ok)
			if ok { // connect to peer
				fmt.Println("Sending to peer...")
				err = udp.Send(conn, &peer, []byte("PING@0"))
				if err != nil {
					fmt.Println(err)
					return // TODO remove it
				}
			} else { // connect to server
				fmt.Println("Sending to server...")
				err = udp.Send(conn, addr, message)
				if err != nil {
					fmt.Println(err)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

	err = udp.Serve(conn, newClientHandler(p))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server is stopped")
	return nil
}

