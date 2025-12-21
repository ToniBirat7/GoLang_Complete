package main

import (
	"fmt"
	"sync"
)

// It is better to pass WaitGroup by pointer or use a global one
var wg sync.WaitGroup

func task(id int, numChan chan int) {
	// 4. This decrements the counter when the function finishes
	defer wg.Done()

	val := <-numChan
	fmt.Printf("Goroutine %d received: %d\n", id, val)
}

func main() {
	messageChan := make(chan int)

	// --- FIRST TASK ---
	wg.Add(1) // 1. Increment counter BEFORE starting goroutine
	go task(1, messageChan)

	messageChan <- 5 // 2. Send value (blocks until G1 is ready to receive)

	// 6. Now Wait() will block because the counter is 2
	wg.Wait()
	fmt.Println("Main: All tasks finished.")
}
