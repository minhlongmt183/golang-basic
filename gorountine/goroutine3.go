package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool)

	fmt.Println("Application start")

	go func() {
		time.Sleep(time.Second)
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}
		done <- true
	}()
	<-done
	fmt.Println("Application end")

}
