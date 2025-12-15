package main

import "fmt"

func main() {
	// Maps
	person := make(map[string]string)

	payload := map[string]int{"person": 40}

	// Add
	person["name"] = "aaaaaaa"

	k, ok := payload["person"]

	fmt.Println(k, ok)
}
