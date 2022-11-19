package main

import (
	"fmt"
	"time"
)

func sendValue(number string, channel chan<- string) {
	for {
		channel <- number
	}
}

func receiveValue(channel <-chan string) {
	for v := range channel {
		fmt.Println(v)
	}
}
func main() {
	channel := make(chan string, 64)
	go sendValue("Hello", channel)
	go sendValue("Edisc", channel)

	go receiveValue(channel)
	time.Sleep(time.Second)
}
