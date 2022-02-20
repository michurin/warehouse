package app

import (
	"bytes"
	"net"
	"time"
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

func taskPing(addr *net.UDPAddr) task {
	return task{
		message:  []byte{labelPing},
		addr:     addr,
		tries:    20,
		interval: 100 * time.Millisecond,
		fin:      false,
	}
}

func taskPong(addr *net.UDPAddr) task {
	return task{
		message:  []byte{labelPong},
		addr:     addr,
		tries:    20,
		interval: 100 * time.Millisecond,
		fin:      false,
	}
}

func taskClose(addr *net.UDPAddr) task {
	return task{
		message:  []byte{labelClose},
		addr:     addr,
		tries:    5,
		interval: 50 * time.Millisecond,
		fin:      true, // it is final task
	}
}

func taskRequestToServer(addr *net.UDPAddr, message []byte) task {
	return task{
		message:  message,
		addr:     addr,
		tries:    -1, // infinite
		interval: 20 * time.Second,
		fin:      false,
	}
}

func taskEexecutor(
	conn Connenction,
	serverAddr *net.UDPAddr,
	serverMessage []byte,
	tq chan task,
	res chan result,
) {
	defaultTask := taskRequestToServer(serverAddr, serverMessage)
	tsk := defaultTask
	ok := true
	for ok {
		_, err := conn.WriteToUDP(tsk.message, tsk.addr)
		if err != nil {
			res <- result{addr: nil, err: err}
			return
		}
		if tsk.tries > 0 {
			tsk.tries--
		}
		if tsk.tries == 0 {
			if tsk.fin {
				res <- result{addr: tsk.addr, err: nil}
				return
			}
			tsk = defaultTask // back to server polling
		}
		select {
		case <-time.After(tsk.interval):
		case tsk, ok = <-tq: // stop looping if channel is closed
		}
	}
}

func serveForever(conn Connenction, tq chan task, res chan result) {
	buff := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buff)
		if err != nil {
			res <- result{addr: nil, err: err}
			break
		}
		ff := bytes.Split(buff[:n], []byte{labelsSeporator}) // TODO we don't need it here
		if len(ff) == 0 {
			continue
		}
		if len(ff[0]) == 0 {
			continue
		}
		switch ff[0][0] {
		case labelPeerInfo:
			peerAddr, err := net.ResolveUDPAddr("udp", string(ff[2]))
			if err != nil {
				// TODO log? stop?
				continue
			}
			tq <- taskPing(peerAddr)
		case labelPing:
			tq <- taskPong(addr)
		case labelPong:
			tq <- taskClose(addr) // task "close" will stop executor after all tries
			return                // stop listening on first pong
		case labelClose:
			close(tq) // stop execution immediately
			res <- result{addr: addr, err: nil}
			return // stop listening on first close
		default:
			// TODO Unexpected data. Log? stop? sleep?
		}
	}
}

func Client(slot, address, remoteAddress string, opt ...Option) (*net.UDPAddr, *net.UDPAddr, error) {
	config := newConfig(opt...)

	message := []byte(slot) // TODO check if slot = 'a'|'b'

	laddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, nil, err
	}
	udpConn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, nil, err
	}
	conn := Connenction(udpConn)
	for _, mw := range config.connMW {
		conn = mw(conn)
	}
	defer conn.Close()

	addr, err := net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		return nil, nil, err
	}

	taskQueue := make(chan task, 8)
	resultChan := make(chan result, 1)

	go taskEexecutor(conn, addr, message, taskQueue, resultChan)
	go serveForever(conn, taskQueue, resultChan)

	res := <-resultChan
	return laddr, res.addr, res.err
}