package main

import (
	"fmt"
)

type Number interface {
	int | string
}

func printSlice[T Number](items []T) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func main() {
	printSlice([]int{1, 2, 3, 4})
	printSlice([]string{"1", "2", "3", "4"})
}