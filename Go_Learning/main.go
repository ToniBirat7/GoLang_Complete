package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	sourceFile, err := os.Open("example.txt")

	if err != nil {
		panic(err)
	}

	destFile, err := os.Create("example2.txt")

	defer sourceFile.Close()
	defer destFile.Close()

	reader := bufio.NewReader(sourceFile)
	writer := bufio.NewWriter(destFile)

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err.Error() != "EOF" {
				panic(err)
			}
			break
		}

		e := writer.WriteByte(b)

		if e != nil {
			panic(e)
		}
	}

	writer.Flush()

	fmt.Print("Completed")
}
