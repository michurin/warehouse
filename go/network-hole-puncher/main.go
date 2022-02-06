package main

/*
args:
  server addr
  client addr server_addr

TODO:
- security
- retries (lost packages)
- timeouts
- refactoring, logging
- smarter contract. check punch
- usage
- readme and instruction
- tests (after refactoring and decomposition)
*/

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Core functions. They return errors
// TODO remove all Prints from this functions, use wrappers for logging

func connect(address string) (*net.UDPConn, error) {
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

func serve(
	conn *net.UDPConn,
	handler func(*net.UDPConn, *net.UDPAddr, []byte) (bool, error),
) error {
	data := make([]byte, 1024)
	for {
		fmt.Println("reading...")
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

func send(conn *net.UDPConn, addr *net.UDPAddr, data []byte) error {
	_, err := conn.WriteToUDP(data, addr)
	return err
}

// Handlers (have to be here)

var (
	prevAddr   = net.UDPAddr{}
	prevAddrMx = new(sync.Mutex) // TODO we really need it?
)

func getAndUpdateAddr(update *net.UDPAddr) *net.UDPAddr {
	prevAddrMx.Lock()
	defer prevAddrMx.Unlock()
	if bytes.Compare(prevAddr.IP, update.IP) == 0 && prevAddr.Port == update.Port {
		return nil
	}
	p := prevAddr
	prevAddr = *update
	if p.IP == nil { // very first interaction
		return nil
	}
	return &p
}

func serverHandler(
	conn *net.UDPConn,
	addr *net.UDPAddr,
	data []byte,
) (bool, error) {
	var payload string
	prev := getAndUpdateAddr(addr)
	if prev == nil {
		payload = "no address"
	} else {
		payload = prev.String()
	}
	err := send(conn, addr, []byte(payload))
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return false, nil
}

func clientHandler(
	conn *net.UDPConn,
	addr *net.UDPAddr,
	data []byte,
) (bool, error) {
	fmt.Printf("THE RESULT: REMOTE GLOBAL ADDRESS: %q\n", data)
	return true, nil // TODO exit only if we obtain address
}

// UI functions. They print logs

func serverMode(address string) {
	conn, err := connect(address)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	err = serve(conn, serverHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func clientMode(address, remoteAddress string) {
	conn, err := connect(address)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	addr, err := net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() { // TODO repeat sending
		fmt.Println("sending...")
		err = send(conn, addr, []byte("MICHURIN"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	err = serve(conn, clientHandler)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Server is stopped")
}

func help() {
	fmt.Println(`Usage: TODO`)
}

func match(op, pat string) bool {
	a := len(op)
	b := len(pat)
	if b < a {
		a = b
	}
	return strings.ToLower(op[:a]) == strings.ToLower(pat[:a])
}

func app(args []string) {
	if len(args) == 0 {
		help()
		return
	}
	op := args[0]
	if match(op, "server") {
		if len(args) != 2 {
			help()
			return
		}
		serverMode(args[1])
		return
	}
	if match(op, "client") {
		if len(args) != 3 {
			help()
			return
		}
		clientMode(args[1], args[2])
		return
	}
}

func main() {
	app(os.Args[1:])
}
