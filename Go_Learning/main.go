package main

import (
	"fmt"
	"sync"
)

type post struct {
	views int
}

var wg sync.WaitGroup

func (p *post) inc(wg *sync.WaitGroup) {
	defer wg.Done()
	p.views += 1
	fmt.Println("Current Views is : ", p.views)
}

func main() {

	myPost := post{
		views: 0,
	}

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go myPost.inc(&wg)
	}

	wg.Wait()
}