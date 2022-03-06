package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/michurin/warehouse/go/network-hole-puncher/app"
)

const (
	role_control = "c"
	role_node_a  = "a"
	role_node_b  = "b"
)

func help() {
	fmt.Fprintf(os.Stderr, `USAGE:
%[1]s role secret local_addr [server_addr]

Roles:
  a — Node A
  b — Node B
  c — Control node for coordination (server)

Examples:
Server mode:
%[1]s c secret :5555
Client mode:
%[1]s a secret :7777 1.2.3.4:5555
%[1]s b secret :7777 1.2.3.4:5555
`, os.Args[0])
}

func helpAndExitIfError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "ERROR:", err.Error())
	help()
	os.Exit(1)
}

func safeIP(ip net.IP) string {
	switch len(ip) {
	case net.IPv4len, net.IPv6len:
		return ip.String()
	default:
		return "n/a"
	}
}

func printResult(laddr, addr *net.UDPAddr) {
	fmt.Println(
		"LADDR/LHOST/LPORT/RADDR/RHOST/RPORT:",
		laddr,
		safeIP(laddr.IP),
		laddr.Port,
		addr,
		safeIP(addr.IP),
		addr.Port)
}

func parseArgsCommon() (role string, secret []byte, laddr string, err error) {
	if len(os.Args) < 4 {
		err = errors.New("args: not enough arguments")
		return
	}
	role = os.Args[1]
	secret = []byte(os.Args[2])
	laddr = os.Args[3]
	return
}

func checkArgs(n int, m string) error {
	if len(os.Args) == n {
		return nil
	}
	return errors.New(m)
}

func main() {
	role, secret, laddr, err := parseArgsCommon()
	helpAndExitIfError(err)

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	logger.SetPrefix(fmt.Sprintf("[%d] [%s] ", os.Getpid(), role))
	opts := app.ConnOption(app.SignMW(secret), app.LogMW(logger))

	switch role {
	case role_control:
		helpAndExitIfError(checkArgs(4, "You have to specify 3 arguments in control (`c`) mode"))
		logger.Print("[INFO] Server started on " + laddr)
		err := app.Server(laddr, opts)
		helpAndExitIfError(err)
		return
	case role_node_a, role_node_b:
		helpAndExitIfError(checkArgs(5, "You have to specify 4 arguments in node (`a` and `b`) mode"))
		raddr := os.Args[4]
		logger.Print("[INFO] Client started on " + laddr + " to server at " + raddr)
		laddr, addr, err := app.Client(role, laddr, raddr, opts) // btw, abstraction leaking (role: arg->payload)
		helpAndExitIfError(err)
		printResult(laddr, addr)
		return
	default:
		help()
	}
}
