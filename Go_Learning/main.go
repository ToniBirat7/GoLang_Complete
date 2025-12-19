package main

import (
	"fmt"
)

func task(id int) {
	defer wg.Done()
	fmt.Println("Doing Task", id)
}

func main() {

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go task(i)
	}

	fmt.Println("Completed")

	wg.Wait()
}
