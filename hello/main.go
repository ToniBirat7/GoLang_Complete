package main

import (
	"fmt"
)

type order struct {
	id     string
	amount float32
	status string
	// createdAt time.Time
}

func (o *order) getStatus() string {
	return o.status
}

func (o *order) changeStatus(status string) {
	o.status = status
}

func NewOrder(id string, amount float32, status string) *order {
	return &order{
		id:     id,
		amount: amount,
		status: status,
	}
}

func main() {
	myOrder2 := NewOrder("1", 200.0, "no")

	myOrder2.changeStatus("Yes")

	fmt.Println(myOrder2.getStatus())
}