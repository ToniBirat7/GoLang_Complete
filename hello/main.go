package main

import (
	"fmt"
)

func closure() func() int {
	count := 0

	return func() int {
		count += 1
		return count
	}
}

func main() {
	fn := closure()

	fmt.Println(fn())
	fmt.Println(fn())
	fmt.Println(fn())
}
