package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("this is a test")

	var mu sync.Mutex

	ch := make(chan int)

	<-ch
	mu.Lock()
	mu.Lock()
}
