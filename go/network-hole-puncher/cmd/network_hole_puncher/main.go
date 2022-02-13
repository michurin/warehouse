package main

import (
	"fmt"
	"os"

	"github.com/michurin/warehouse/go/network-hole-puncher/app"
)

func help(a string) { // TODO
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
`, a)
}

func main() {
	if len(os.Args[1:]) < 2 {
		help(os.Args[0])
		return
	}
	switch os.Args[1] {
	case "c":
		if len(os.Args) == 3 {
			app.Server(os.Args[2])
			return
		}
	case "a", "b":
		if len(os.Args) == 4 {
			addr, err := app.Client(os.Args[1], os.Args[2], os.Args[3])
			if err != nil {
				help((os.Args[0]))
				return
			}
			fmt.Println("RESULT:", addr)
			fmt.Printf("remote %s %d udp\n", addr.IP, addr.Port)
			return
		}
	}
	fmt.Println("Invalid arguments")
	help(os.Args[0])
}
