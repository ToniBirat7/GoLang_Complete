package main

import "fmt"

type post struct {
	views int
}

func (p *post) inc() {
	p.views += 1
	fmt.Println("Current Views is : ", p.views)
}

func callInc(myPost post) {
	myPost.inc()
}

func main() {

	myPost := post{
		views: 0,
	}

	for i := 0; i <= 10; i++ {
		go callInc(myPost)
	}
}
