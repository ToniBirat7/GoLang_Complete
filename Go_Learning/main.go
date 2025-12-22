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
}

func main() {

	myPost := post{
		views: 0,
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go myPost.inc(&wg)
	}

	wg.Wait()

	fmt.Println("Final Sum : ", myPost.views)
}
