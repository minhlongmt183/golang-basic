package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	fmt.Println("Application start")

	wg.Add(1) // 1: goi 1 routine
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("Goroutines: ", i)
		}

		wg.Done()
	}()
	wg.Wait()
	fmt.Println("Application end")

}
