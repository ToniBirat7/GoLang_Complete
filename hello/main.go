package main

import (
	"fmt"
	"time"
)

func blockingTask(id int) {
	fmt.Printf("Goroutine %d: Starting blocking task\n", id)
	time.Sleep(2 * time.Second) // Simulating a blocking operation
	fmt.Printf("Goroutine %d: Completed blocking task\n", id)
}

func main() {
	for i := 1; i <= 3; i++ {
		go blockingTask(i)
	}
	time.Sleep(5 * time.Second) // Sleep to allow Goroutines to complete
}
