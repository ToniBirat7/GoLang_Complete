package main

import (
	"fmt"
	"sync"
)

func task(id int, w *sync.WaitGroup) {
	fmt.Println("Doing Task", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go task(i, &wg)
	}
}