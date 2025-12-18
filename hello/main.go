package main

import "fmt"

func printSlice(items []int) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func main() {
	printSlice([]int{1, 2, 3, 4})
}