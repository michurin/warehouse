package app

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/michurin/warehouse/go/network-hole-puncher/internal/udp"
)

type task struct {
	message  []byte
	addr     *net.UDPAddr
	tries    int
	interval time.Duration
	fin      bool
}

type result struct {
	addr *net.UDPAddr
	err  error
}

func newClientHandler(tq chan task, res chan result) udp.Handler {
	return func(
		conn *net.UDPConn,
		addr *net.UDPAddr,
		data []byte,
	) {
		fmt.Printf("DATA: %q\n", data)
		ff := bytes.Split(data, []byte{'@'})
		fmt.Println(ff)
		if len(ff) == 0 { // TODO we have to have this check; HOWEVER server has to return correct signature
			return
		}
		if len(ff[0]) == 0 {
			return
		}
		// TODO |
		// TODO | ff[0][1] is ugly! The message format has to be reinvented
		// TODO |
		if ff[0][1] == 'E' { // (PEER) we receive peer information
			peerAddr, err := net.ResolveUDPAddr("udp", string(ff[2]))
			if err != nil {
				// TODO log?
				return
			}
			tq <- task{
				message:  []byte("PING"),
				addr:     peerAddr,
				tries:    10, // increase?
				interval: time.Second,
				fin:      false,
			}
			return
		}
		if ff[0][1] == 'I' { // (PING) from host
			tq <- task{
				message:  []byte("PONG"),
				addr:     addr,
				tries:    10,
				interval: 100 * time.Millisecond,
				fin:      false, // TODO maybe true?
			}
			return
		}
		if ff[0][1] == 'O' { // (PONG) from host
			tq <- task{
				message:  []byte("CLOSE"),
				addr:     addr,
				tries:    5,
				interval: 100 * time.Millisecond,
				fin:      true, // it is final task
			}
			return
		}
		if ff[0][1] == 'L' { // (CLOSE) from host (fast done)
			fmt.Println("FIN BY CLOSE")
			res <- result{addr: addr, err: nil}
			return
		}
		// TODO unexpected data; error that has to be ignore. It evaporates when protocol gets more strict and secure
	}
}

func taskEexecutor(conn *net.UDPConn, serverAddr *net.UDPAddr, serverMessage []byte, tq chan task, res chan result) { // TODO repeat sending
	defaultTask := task{
		message:  serverMessage, // request to server
		addr:     serverAddr,    // address of server
		tries:    -1,            // infinite
		interval: 10 * time.Second,
	}
	tsk := defaultTask
	for {
		// execute task
		fmt.Println("Execute task:", string(tsk.message), tsk)
		err := udp.Send(conn, tsk.addr, tsk.message)
		if err != nil {
			res <- result{addr: nil, err: err}
			return
		}
		if tsk.tries > 0 {
			tsk.tries--
		}
		if tsk.tries == 0 {
			if tsk.fin {
				fmt.Println("FIN BY FIN")
				res <- result{addr: tsk.addr, err: nil}
				return
			}
			tsk = defaultTask // back to server polling
		}
		select {
		case <-time.After(tsk.interval):
			fmt.Println("timeout")
		case tsk = <-tq:
			fmt.Println("new task", tsk.message)
		}
	}
}

func serveForever(conn *net.UDPConn, h udp.Handler, res chan result) {
	err := udp.Serve(conn, h)
	res <- result{addr: nil, err: err}
}

func Client(slot, address, remoteAddress string) (*net.UDPAddr, error) {
	message := []byte(slot) // TODO check if slot = 'a'|'b'

	conn, err := udp.Connect(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	addr, err := net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		return nil, err
	}

	taskQueue := make(chan task, 8)
	resultChan := make(chan result, 1)

	go taskEexecutor(conn, addr, message, taskQueue, resultChan)
	go serveForever(conn, newClientHandler(taskQueue, resultChan), resultChan)

	res := <-resultChan
	return res.addr, res.err
}
