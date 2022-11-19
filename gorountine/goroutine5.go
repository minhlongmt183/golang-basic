package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan string)

	fmt.Println("Application start")

	done <- "Done"
	fmt.Println(<-done)

	fmt.Println("Application end")
	time.Sleep(time.Second)
}
