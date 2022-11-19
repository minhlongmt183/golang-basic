package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	connType = "tcp"
)

func RevShellListener(port string, done chan<- bool) {
	fmt.Printf("Starting listener...\n")
	port = ":" + port
	conn, err := net.Listen("tcp4", port)
	if err != nil {
		fmt.Println("listening error: ", err)
		return
	}
	defer func() {
		conn.Close()
		done <- true
	}()

	rand.Seed(time.Now().Unix())

	for {
		c, err := conn.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	//connection closed
	defer conn.Close()

	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	reader := bufio.NewReader(os.Stdin)

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Println("close connection...")
			return
		}
	}()
	for {
		toExec, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("read error")
			break
		}

		temp := strings.TrimSpace(string(toExec))
		//fmt.Printf("Executing: %s\n", temp)
		_, err = fmt.Fprintf(conn, "%s\n", temp)
		if err != nil {
			fmt.Println("Get error: " + err.Error())
			continue
		}

	}

}

func main() {
	arguments := os.Args
	done := make(chan bool)
	if len(arguments) == 1 {
		fmt.Errorf("Usage: ./reverse-tcp-server <port_number>")
		return
	}

	fmt.Printf("Starting\n")

	go RevShellListener(arguments[1], done)
	<-done
	fmt.Println("Finished")
}
