package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func swap(a int, b int) (int, int) {
	return b, a
}

func main() {
	sum := add(2, 3)
	fmt.Println(sum)

	x, y := swap(2, 3)
	fmt.Println(x, y)
}
