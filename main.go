package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	ip := "127.0.0.1"        // CHANGE THIS (IP)
	port := [2]int{0, 65535} // CHANGE THIS (port range) default (0-65535)

	fmt.Printf("Scanning host %v, ports %v-%v\n", ip, port[0], port[1])

	var wg sync.WaitGroup
	startedScanAt := time.Now()
	openPortsCount := 0

	for i := port[0]; i <= port[1]; i++ {
		wg.Add(1)
		go scanPort(ip, i, &openPortsCount, &wg)
	}

	wg.Wait()
	fmt.Println("scanning completed in", time.Since(startedScanAt))
	fmt.Println("open ports:", openPortsCount)
}

func scanPort(ip string, port int, openPortsCount *int, wg *sync.WaitGroup) bool {
	_, err := net.Dial("tcp", fmt.Sprintf("%v:%v", ip, port))
	wg.Done()

	if err == nil {
		fmt.Printf("%v port is open\n", port)
		*openPortsCount++
		return true
	}

	return false
}
