package main

import (
	"fmt"
	"sync"
)

// It is better to pass WaitGroup by pointer or use a global one
var wg sync.WaitGroup

func receiver(numChan chan int) {
	// 4. This decrements the counter when the function finishes
	defer wg.Done()

	val := <-numChan
	fmt.Printf("Goroutine received: %d\n", val)
}

func sender(value int, messageChan chan int) {
	messageChan <- 5
}

func main() {
	messageChan := make(chan int)

	wg.Add(1)
	go receiver(messageChan)

	wg.Add(1)
	go sender(5, messageChan)

	wg.Wait()
	fmt.Println("Main: All tasks finished.")
}
