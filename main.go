package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage("No monit target was provided.", "")
		os.Exit(1)
	}
	monit(args)
}

func usage(errors ...string) {
	usage := "Usage: " +
		os.Args[0] +
		"<hostname | ip address> [hostname 2 | ip address 2] [...]" +
		"\n"
	for e := range errors {
		usage += "\n" + errors[e]
	}
	fmt.Fprintf(os.Stderr, usage)
}

func monit(addresses []string) {
	var res []chan string
	for n := range addresses {
		res = append(res, make(chan string))
		go checkSrv(addresses[n], res[n])
	}
	prt := <-res[0]
	fmt.Println(prt)
}

func checkSrv(addr string, ret chan string) {
	res, _ := exec.Command("ping", addr, "-c 3").Output()
	if strings.Contains(string(res), "Destination Host Unreachable") ||
		strings.Contains(string(res), "100% packet loss") {
		ret <- "OFFLINE"
	} else {
		pingRows := strings.Split(string(res), "\n")
		pingRow := pingRows[len(pingRows)-2]
		pingSlc := strings.Split(pingRow, "/")
		ret <- pingSlc[len(pingSlc)-3]
	}
}
