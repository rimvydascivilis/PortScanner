package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var portRange = flag.String("p", "0:65535", "Port range, syntax: -p <start>:<end>")

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	startedScanAt := time.Now()
	openPortsCount := 0
	host, portRangeStart, portRangeEnd := getFlags()

	wg.Add(portRangeEnd - portRangeStart + 1)

	fmt.Printf("Scanning host %v, ports %v-%v\n", host, portRangeStart, portRangeEnd)

	for i := portRangeStart; i <= portRangeEnd; i++ {
		go tcpScan(host, i, &openPortsCount, &wg)
	}

	wg.Wait()
	fmt.Println("Open ports:", openPortsCount)
	fmt.Println("Scanning completed in", time.Since(startedScanAt))
}

func tcpScan(host string, port int, openPortsCount *int, wg *sync.WaitGroup) {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))

	if err == nil {
		fmt.Printf("[+] %v port is open\n", port)
		*openPortsCount++
	}
	wg.Done()
}

func getFlags() (string, int, int) {
	host := flag.CommandLine.Arg(0)
	portRangeSlice := strings.Split(*portRange, ":")
	portRangeStart, err := strconv.Atoi(portRangeSlice[0])
	if err != nil {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
	portRangeEnd, err := strconv.Atoi(portRangeSlice[1])
	if err != nil {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
	verifyFlags(host, portRangeSlice, portRangeStart, portRangeEnd)

	return host, portRangeStart, portRangeEnd
}

func verifyFlags(host string, portRangeSlice []string, portRangeStart int, portRangeEnd int) {
	if len(portRangeSlice) != 2 {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
	if host == "" {
		fmt.Println("You must Specify a host")
		os.Exit(1)
	}
	if portRangeEnd < portRangeStart {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
}
