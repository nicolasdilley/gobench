package main

import (
	"sync"

)

type processorListener struct {
	lock sync.RWMutex
	cond sync.Cond

	pendingNotifications []interface{}
}

func (p *processorListener) add(notification interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pendingNotifications = append(p.pendingNotifications, notification)
	p.cond.Broadcast()
}

func (p *processorListener) pop(stopCh <-chan struct{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for {
		for len(p.pendingNotifications) == 0 {
			select {
			case <-stopCh:
				return
			default:
			}
			p.cond.Wait()
		}
		select { // block here
		case <-stopCh:
			return
		}
	}
}

func newProcessListener() *processorListener {
	ret := &processorListener{
		pendingNotifications: []interface{}{},
	}
	ret.cond.L = &ret.lock
	return ret
}
func main() {
	pl := newProcessListener()
	stopCh := make(chan struct{})
	defer close(stopCh)
	pl.add(1)
	go pl.pop(stopCh)

	resultCh := make(chan struct{})
	go func() {
		pl.lock.Lock() // block here
		close(resultCh)
	}()
	<-resultCh
	pl.lock.Unlock()
}



