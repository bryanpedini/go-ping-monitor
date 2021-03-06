package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Ping is the join between a pingable host and its immediate result,
// it is updated over time.
type Ping struct {
	host   string
	result string
}

func main() {
	args := os.Args[1:]
	// checking for at least one provided target
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
	pings := make(chan Ping, len(addresses))
	var wg sync.WaitGroup
	for _, addr := range addresses {
		wg.Add(1)
		go checkSrv(addr, pings, &wg)
	}
	wg.Wait()
	close(pings)
	for res := range pings {
		fmt.Printf("%s %s", res.host, res.result)
	}
	fmt.Println()
}

func checkSrv(addr string, ret chan Ping, wg *sync.WaitGroup) {
	defer wg.Done()
	p := Ping{host: addr}
	res, err := exec.Command("ping", addr, "-c 3").Output()
	if err != nil || strings.Contains(string(res), "Destination Host Unreachable") ||
		strings.Contains(string(res), "100% packet loss") {
		p.result = "OFFLINE"
	} else {
		pingRows := strings.Split(string(res), "\n")
		pingRow := pingRows[len(pingRows)-2]
		pingSlc := strings.Split(pingRow, "/")
		p.result = pingSlc[len(pingSlc)-3]
	}
	ret <- p
}
