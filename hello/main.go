package main

import (
	"fmt"
)

type PaymentMethod interface {
	ProcessPayment(amount float32) string
}

type CreditCard struct {
	cardNumber string
}

func (cc CreditCard) ProcessPayment(amount float32) string {
	return fmt.Sprintf("Processing credit card payment of $%.2f using card number %s", amount, cc.cardNumber)
}

type PayPal struct {
	email string
}

func (cc PayPal) ProcessPayment(amount float32) string {
	return fmt.Sprintf("Processing PayPal card payment of $%.2f using card number %s", amount, cc.email)
}

func processPayment(method PaymentMethod, amount float32) {
	result := method.ProcessPayment(amount)
	fmt.Println(result)
}

func main() {
	cc := CreditCard{cardNumber: "1234-5678"}
	pp := PayPal{email: "hi@"}
	processPayment(cc, 100.0)
	processPayment(pp, 50.0)
}