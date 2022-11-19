package main

import (
	"fmt"
)

func main() {
	fmt.Println("Application start")

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}
	}()

	fmt.Println("Application end")
	//time.Sleep(time.Second)
}
