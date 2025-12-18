package main

import (
	"fmt"
)

type Number interface {
	int | string
}

type Pair[T any, U any] struct {
	first  T
	second U
}

func printSlice[T Number](items []T) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func main() {
	pair := Pair[string, bool]{
		first:  "birat",
		second: true,
	}
	fmt.Println(pair.first)
	printSlice([]int{1, 2, 3, 4})
	printSlice([]string{"1", "2", "3", "4"})
}
