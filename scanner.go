// Simple Concurrent Port Scanner in Go
// Author: MrZefDev
// Educational use only.

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var timeout = 800 * time.Millisecond

func scanPort(host string, port int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return
	}

	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)
	banner, _ := reader.ReadString('\n')

	conn.Close()

	if banner != "" {
		fmt.Printf("[+] %d OPEN -> %s", port, banner)
	} else {
		fmt.Printf("[+] %d OPEN\n", port)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run scanner.go <host> <startPort> <endPort>")
		return
	}

	host := os.Args[1]
	start, _ := strconv.Atoi(os.Args[2])
	end, _ := strconv.Atoi(os.Args[3])

	var wg sync.WaitGroup
	sem := make(chan struct{}, 200)

	fmt.Println("Target:", host)
	fmt.Printf("Scanning ports %d-%d...\n", start, end)

	startTime := time.Now()

	for port := start; port <= end; port++ {
		wg.Add(1)
		go scanPort(host, port, &wg, sem)
	}

	wg.Wait()

	fmt.Println("\nScan finished.")
	fmt.Println("Elapsed:", time.Since(startTime))
}
