package main

import (
	"fmt"
	"os"

	"github.com/michurin/warehouse/go/network-hole-puncher/app"
)

func help(a string) { // TODO
	fmt.Printf(`USAGE
Server mode:
%[1]s s :5555
Client mode:
%[1]s c 1 :7777 1.2.3.4:5555
`, a)
}

func main() {
	if len(os.Args[1:]) < 2 {
		help(os.Args[0])
		return
	}
	if os.Args[1] == "s" && len(os.Args) == 3 {
		app.Server(os.Args[2])
		return
	}
	if os.Args[1] == "c" && len(os.Args) == 5 {
		addr, err := app.Client(os.Args[2], os.Args[3], os.Args[4])
		if err != nil {
			help((os.Args[0]))
			return
		}
		fmt.Println("RESULT:", addr)
		fmt.Printf("remote %s %d udp\n", addr.IP, addr.Port)
		return
	}
	help(os.Args[0])
}
