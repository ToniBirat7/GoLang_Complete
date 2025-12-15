package main

import (
	"fmt"
)

func changeNum(num *int) {
	fmt.Println("In ChangeNum : ", *num)
	*num = 5 // Derefencing
}

func main() {
	num := 1

	changeNum(&num)

	fmt.Println("After ChangeNum : ", num)
}
