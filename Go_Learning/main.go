package main

import (
	"fmt"
	"sync"
)

func task(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Doing Task", id)
}

var wg sync.WaitGroup

func main() {

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go task(i, &wg)
	}

	fmt.Println("Completed")

	wg.Wait()
}
