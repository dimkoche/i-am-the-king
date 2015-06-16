package main

import "net"
import "fmt"
import "bufio"
import "os"
import "time"

var T time.Duration
var PRIORITY_MAP map[string]string

func main() {
	initGlobals()

	priority, ok := getPriority()
	if !ok {
		os.Exit(1)
	}

	server := &Server{
		priority: priority,
		state:    "init",
		joins:    make(chan net.Conn),
	}
	server.Start()

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Printf("CMD: %s", text)
	}
}

func initGlobals() {
	PRIORITY_MAP = map[string]string{
		"1": "localhost:8081",
		"2": "localhost:8082",
		"3": "localhost:8083",
		"4": "localhost:8084",
		"5": "localhost:8085",
		"6": "localhost:8086",
		"7": "localhost:8087",
	}

	T = 3 * time.Second
}

func getPriority() (string, bool) {
	if len(os.Args) < 2 {
		fmt.Println("Print help")
		return "", false
	}

	priority := os.Args[1]
	_, ok := PRIORITY_MAP[priority]
	if !ok {
		fmt.Println("Priority not found:", priority)
		return "", false
	}

	fmt.Println("Launching server with priority:", priority, ". Address:", PRIORITY_MAP[priority])
	return priority, true
}
