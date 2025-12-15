package main

import "fmt"

func applyOperation(a int, b int, operation func(int, int) int) int {
	return operation(a, b)
}

func add(a int, b int) int {
	return a + b
}

func main() {
	result := applyOperation(2, 3, add)
	fmt.Println(result)
}
