package main

import (
	"fmt"
	"os"
)

func main() {
	dir, err := os.Open(".")
	if err != nil {
		panic(err)
	}

	defer dir.Close()

	fileInfo, err := dir.ReadDir(0)

	for _, fi := range fileInfo {
		fmt.Println(fi.Name())
	}

}
