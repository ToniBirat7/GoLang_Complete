package main

import (
	"fmt"
)

type order struct {
	id     string
	amount float32
	status string
	// createdAt time.Time // Nanosecond precision
}

func main() {
	myOrder := order{
		id:     "123",
		amount: 1200.00,
		status: "received",
		// createdAt: time.Now(),
	}

	fmt.Println("Order Struct", myOrder)
}