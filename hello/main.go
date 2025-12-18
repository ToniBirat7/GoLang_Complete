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

func main() {
	myOrder := order{
		id:     "1",
		amount: 200.0,
		status: "no",
	}
	// fmt.Println(myOrder)

	myOrder.changeStatus("Yes")

	fmt.Println(myOrder.getStatus())
}
