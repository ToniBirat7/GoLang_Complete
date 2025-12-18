package main

import (
	"fmt"
	"time"
)

type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time
}

func main() {
	myOrder := order{
		id:     "1",
		amount: 200.0,
		status: "no",
	}

	fmt.Println(myOrder)
}
