package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func task(numChan chan int) {
	fmt.Println("Received Value is : ", <-numChan)
}

func main() {

	messageChan := make(chan int)

	go task(messageChan)

	messageChan <- 5

	wg.Wait()
}
