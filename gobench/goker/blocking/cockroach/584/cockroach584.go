package main

import (
	"sync"

)

type Gossip struct {
	mu     sync.Mutex
	closed bool
}

func (g *Gossip) bootstrap() {
	for {
		g.mu.Lock()
		if g.closed {
			/// Missing g.mu.Unlock
			break
		}
		g.mu.Unlock()
		break
	}
}

func (g *Gossip) manage() {
	for {
		g.mu.Lock()
		if g.closed {
			/// Missing g.mu.Unlock
			break
		}
		g.mu.Unlock()
		break
	}
}
func main() {
	g := &Gossip{
		closed: true,
	}
	go func() {
		g.bootstrap()
		g.manage()
	}()
}



