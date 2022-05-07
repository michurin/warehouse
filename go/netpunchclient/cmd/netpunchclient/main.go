package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/michurin/netpunch/netpunchlib"
)

var logger = log.Default()

func checkErr(err error) {
	if err != nil {
		logger.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 8 {
		logger.Fatalf("Usage: %s a public-control-node.com:9999 secret openvpn 192.168.2.1 192.168.2.2 /etc/secret.key", filepath.Base(os.Args[0]))
	}

	laddr, addr, err := netpunchlib.Client(
		context.Background(),
		os.Args[1], // arg #1: role: 'a' or 'b'
		"0.0.0.0:"+strconv.Itoa(rand.Intn(8000)+2000),
		os.Args[2], // arg #2: public control node address: example.com:9999
		netpunchlib.ConnOption(
			netpunchlib.SigningMiddleware([]byte(os.Args[3])), // arg #3: secret
			netpunchlib.LoggingMiddleware(logger)))
	checkErr(err)

	executable, err := exec.LookPath(os.Args[4]) // arg #4: openvpn bin
	checkErr(err)
	cli := []string{
		executable, // BTW, you are able to add sudo in front of executable
		"--remote", addr.IP.String(), "--rport", strconv.Itoa(addr.Port),
		"--local", laddr.IP.String(), "--lport", strconv.Itoa(laddr.Port),
		"--proto", "udp", "--dev", "tun",
		"--ifconfig", os.Args[5], os.Args[6], // args #5 and #6: local and remote VPN addresses
		"--auth-nocache", "--secret", os.Args[7], "--auth", "SHA256", "--cipher", "AES-256-CBC", // arg #7: openvpn secret
		"--ping", "10", "--ping-exit", "40",
		"--verb", "3",
	}
	logger.Println("Going to execute:", strings.Join(cli, " "))
	err = syscall.Exec(cli[0], cli, os.Environ())
	checkErr(err)
}
