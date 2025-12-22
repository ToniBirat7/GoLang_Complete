package main

import (
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("example.txt")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	buf := make([]byte, 100)

	d, err := f.Read(buf)

	if err != nil {
		panic(err)
	}

	fmt.Println("Data : ", d, string(buf))
}
