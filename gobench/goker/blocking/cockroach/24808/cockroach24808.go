package main

import (
	"context"
	"sync"

)

type Compactor struct {
	ch chan struct{}
}

type Stopper struct {
	stop    sync.WaitGroup
	stopper chan struct{}
}

func (s *Stopper) Stop() {
	close(s.stopper)
}

func (c *Compactor) Start(ctx context.Context, stopper *Stopper) {
	c.ch <- struct{}{}

	stopper.stop.Add(1)
	go func() {
		defer stopper.stop.Done()
		for {
			select {
			case <-stopper.stopper:
				return
			case <-c.ch:
			}
		}
	}()
}

func main() {
	stopper := &Stopper{
		stopper: make(chan struct{}),
	}
	defer stopper.Stop()

	compactor := &Compactor{ch: make(chan struct{}, 1)}
	compactor.ch <- struct{}{}

	compactor.Start(context.Background(), stopper)
}



