package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/michurin/warehouse/go/network-hole-puncher/app"
)

func help(err error) { // TODO
	if err != nil {
		fmt.Println("\nERROR:", err.Error())
	}
	fmt.Printf(`USAGE:
%[1]s role secret local_addr [server_addr]

Roles:
  a — Node A
  b — Node B
  c — Control node for coordination (server)

Examples:
Server mode:
%[1]s c :5555
Client mode:
%[1]s a :7777 1.2.3.4:5555
%[1]s b :7777 1.2.3.4:5555
`, os.Args[0])
}

func safeIP(ip net.IP) string {
	switch len(ip) {
	case net.IPv4len, net.IPv6len:
		return ip.String()
	}
	return "n/a"
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

func parseArgs() (role string, secret []byte, laddr, raddr string) {
	l := len(os.Args)
	if l < 4 {
		return
	}
	role = os.Args[1]
	secret = []byte(os.Args[2])
	laddr = os.Args[3]
	if l < 5 {
		return
	}
	raddr = os.Args[4]
	return
}

func main() {
	role, secret, laddr, raddr := parseArgs()

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	opts := app.ConnOption(app.SignMW(secret), app.LogMW(logger))

	switch role {
	case "c":
		logger.Print("[INFO] Server started on " + laddr)
		err := app.Server(laddr, opts)
		if err != nil {
			help(err)
			return
		}
		return
	case "a", "b":
		logger.Print("[INFO] Client started on " + laddr)
		laddr, addr, err := app.Client(role, laddr, raddr, opts)
		if err != nil {
			help(err)
			return
		}
		printResult(laddr, addr)
		return
	}

	help(nil)
}
