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
		tries:    10, // increase?
		interval: time.Second,
		fin:      false,
	}
}

func taskPong(addr *net.UDPAddr) task {
	return task{
		message:  []byte{labelPong},
		addr:     addr,
		tries:    10,
		interval: 100 * time.Millisecond,
		fin:      false, // TODO maybe true?
	}
}

func taskClose(addr *net.UDPAddr) task {
	return task{
		message:  []byte{labelClose},
		addr:     addr,
		tries:    5,
		interval: 100 * time.Millisecond,
		fin:      true, // it is final task
	}
}

func taskRequestToServer(addr *net.UDPAddr, message []byte) task {
	return task{
		message:  message,
		addr:     addr,
		tries:    -1, // infinite
		interval: 10 * time.Second,
	}
}

func taskEexecutor(conn Connenction, serverAddr *net.UDPAddr, serverMessage []byte, tq chan task, res chan result) {
	defaultTask := taskRequestToServer(serverAddr, serverMessage)
	tsk := defaultTask
	for {
		// execute task
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
		case tsk = <-tq:
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
		if len(ff) == 0 {                                    // TODO we have to have this check; HOWEVER server has to return correct signature
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
			tq <- taskClose(addr)
		case labelClose:
			res <- result{addr: addr, err: nil}
		default:
			// TODO Unexpected data. Log? stop?
		}
	}
}

func Client(slot, address, remoteAddress string, opt ...Option) (*net.UDPAddr, error) {
	config := newConfig(opt...)

	message := []byte(slot) // TODO check if slot = 'a'|'b'

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	conn := Connenction(udpConn)
	for _, mw := range config.connMW {
		conn = mw(conn)
	}
	defer conn.Close()

	addr, err = net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		return nil, err
	}

	taskQueue := make(chan task, 8)
	resultChan := make(chan result, 1)

	go taskEexecutor(conn, addr, message, taskQueue, resultChan)
	go serveForever(conn, taskQueue, resultChan)

	res := <-resultChan
	return res.addr, res.err
}
