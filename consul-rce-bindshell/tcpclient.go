package main

/* A simple tcp client. This is nowhere near functional or complete, I am simply keeping it here for now.*/

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {

		fmt.Println("USAGE: ./tcpclient <command_>")
		return
	}
	var tcpPort string = "5555"
	connection, err := net.Dial("tcp", "127.0.0.1:"+tcpPort) //connect to the socket
	if err != nil {
		fmt.Println("An error occurred trying to connect to the target:", err)
	}
	defer connection.Close()

	text := os.Args[1]

	//write input to tcp socket
	fmt.Fprintf(connection, text+"\n")

	//formats and writes to a given io.Writer object, in this case the connection
	recvBuf := make([]byte, 1024)
	_, err = connection.Read(recvBuf[:]) // recv data
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			recvBuf = []byte("")
			// time out
		} else {
			log.Println("read error:", err)
			// some error else, do something else, for example create new conn
		}
	}
	//fmt.Print("Message from server: " + message)
	fmt.Println(string(recvBuf))

}
