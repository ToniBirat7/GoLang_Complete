package main

import "fmt"

func task(id int) {
	fmt.Println("Doing Something", id)
}

func main() {
	for i := 0; i <= 10; i++ {
		go task(i)
	}
}
