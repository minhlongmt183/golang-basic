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

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Errorf("Usage: ./reverse-tcp-server <port_number>")
		return
	}
	fmt.Printf("Starting\n")
	PORT := ":" + arguments[1]
	conn, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
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
			return
		}
	}()

	for {
		fmt.Println("Reading...")
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
