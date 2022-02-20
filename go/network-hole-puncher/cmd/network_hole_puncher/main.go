package main

import (
	"fmt"
	"log"
	"os"

	"github.com/michurin/warehouse/go/network-hole-puncher/app"
)

func help(err error) { // TODO
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	fmt.Printf(`USAGE:
%[1]s role local_addr [server_addr]

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

func main() {
	if len(os.Args[1:]) < 2 {
		help(nil)
		return
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	opts := app.ConnOption(app.SignMW([]byte("SECRET")), app.LogMW(logger))

	switch os.Args[1] {
	case "c":
		if len(os.Args) == 3 {
			app.Server(os.Args[2], opts)
			return
		}
	case "a", "b":
		if len(os.Args) == 4 {
			addr, err := app.Client(os.Args[1], os.Args[2], os.Args[3], opts)
			if err != nil {
				help(err)
				return
			}
			fmt.Println("RESULT:", addr)
			fmt.Printf("remote %s %d udp\n", addr.IP, addr.Port)
			return
		}
	}

	help(nil)
}
