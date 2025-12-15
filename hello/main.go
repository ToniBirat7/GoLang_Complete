package main

import (
	"fmt"
	"time"
)

type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time // Nanosecond precision
}

func (o *order) updateStatus(newStatus string) {
	o.status = newStatus
}

func (o order) getAmount() float32 {
	return o.amount
}

func main() {
	myOrder := order{
		id:        "123",
		amount:    1200.00,
		status:    "received",
		createdAt: time.Now(),
	}

	fmt.Println("Initial Order Status:", myOrder.status)
	myOrder.updateStatus("shipped")
	fmt.Println("Updated Order Status:", myOrder.status)

	fmt.Println("Amount is : ", myOrder.getAmount())
}
