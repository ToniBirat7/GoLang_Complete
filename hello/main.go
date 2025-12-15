package main

import (
	"fmt"
)

func changeNum(num int) {
	num = 5
	fmt.Println("In ChangeNum : ", num)
}

func main() {
	num := 1

	changeNum(num)

	fmt.Println("After ChangeNum : ", num)
}
