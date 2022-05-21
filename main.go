package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var host = flag.String("h", "localhost", "Host to scan")
var portRange = flag.String("p", "0:65535", "Port range, syntax: -p <start>:<end>")
var nmapEnabled = flag.Bool("n", false, "Enable nmap scanning, you need to have nmap installed")

func main() {
	var openPorts []int
	host, portRangeStart, portRangeEnd := getFlags()

	findOpenTcpPorts(host, portRangeStart, portRangeEnd, &openPorts)

	if *nmapEnabled {
		if len(openPorts) == 0 {
			fmt.Println("Nmap scan not started, no open ports found")
			os.Exit(0)
		}
		startNmapScan(host, &openPorts)
	}
}

func findOpenTcpPorts(host string, portRangeStart int, portRangeEnd int, openPorts *[]int) {
	startedScanAt := time.Now()
	var wg sync.WaitGroup
	wg.Add(portRangeEnd - portRangeStart + 1)

	fmt.Printf("Scanning host %v, ports %v-%v\n", host, portRangeStart, portRangeEnd)

	for i := portRangeStart; i <= portRangeEnd; i++ {
		go tcpScan(host, i, openPorts, &wg)
	}

	wg.Wait()
	fmt.Println("Open ports:", len(*openPorts))
	fmt.Println("Scanning completed in", time.Since(startedScanAt))
}

func tcpScan(host string, port int, openPorts *[]int, wg *sync.WaitGroup) {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))

	if err == nil {
		fmt.Printf("[+] %v port is open\n", port)
		*openPorts = append(*openPorts, port)
	}
	wg.Done()
}

func startNmapScan(host string, openPorts *[]int) {
	var nmapArgs []string

	nmapArgs = append(nmapArgs, "-p"+arrayToString(openPorts, ","))
	nmapArgs = append(nmapArgs, "-A")
	nmapArgs = append(nmapArgs, host)

	fmt.Printf("Starting nmap scan 'nmap %v'\n", strings.Join(nmapArgs, " "))

	cmd := exec.Command("nmap", nmapArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running nmap:", err)
	}
}

func getFlags() (string, int, int) {
	flag.Parse()
	portRangeSlice := strings.Split(*portRange, ":")
	portRangeStart, err := strconv.Atoi(portRangeSlice[0])
	host := *host

	if err != nil {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
	portRangeEnd, err := strconv.Atoi(portRangeSlice[1])
	if err != nil {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}

	verifyFlags(portRangeSlice, portRangeStart, portRangeEnd)
	return host, portRangeStart, portRangeEnd
}

func verifyFlags(portRangeSlice []string, portRangeStart int, portRangeEnd int) {
	if len(portRangeSlice) != 2 {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
	if portRangeEnd < portRangeStart {
		fmt.Println("Invalid port range")
		os.Exit(1)
	}
}

func arrayToString(a *[]int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(*a), " ", delim, -1), "[]")
}
