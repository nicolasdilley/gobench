/*
 * Project: grpc-go
 * Issue or PR  : https://github.com/grpc/grpc-go/pull/1424
 * Buggy version: 39c8c3866d926d95e11c03508bf83d00f2963f91
 * fix commit-id: 64bd0b04a7bb1982078bae6a2ab34c226125fbc1
 * Flaky: 100/100
 * Description:
 *   The parent function could return without draining the done channel.
 */
package main

import (
	"sync"

)

type roundRobin struct {
	mu     sync.Mutex
	addrCh chan bool
}

type addrConn struct {
	mu sync.Mutex
}

func (ac *addrConn) tearDown() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
}

type dialOptions struct {
	balancer *roundRobin
}

type ClientConn struct {
	dopts dialOptions
	conns []*addrConn
}

func (cc *ClientConn) lbWatcher(doneChan chan bool) {
	for addr := range cc.dopts.balancer.addrCh {
		if addr {
			// nop, make compiler happy
		}
		var (
			/// add []Address is empty
			del []*addrConn
		)
		for _, a := range cc.conns {
			del = append(del, a)
		}
		for _, c := range del {
			c.tearDown()
		}
		/// Without close doneChan
		/// FIX: defer close(doneChan)
	}
}

func DialContext() {
	cc := &ClientConn{
		dopts: dialOptions{
			balancer: &roundRobin{addrCh: make(chan bool)},
		},
	}
	waitC := make(chan error, 1)
	go func() { // G2
		defer close(waitC)
		ch := cc.dopts.balancer.addrCh
		if ch != nil {
			doneChan := make(chan bool)
			go cc.lbWatcher(doneChan) // G3
			<-doneChan                /// Block here
		}
	}()
	/// close addrCh
	close(cc.dopts.balancer.addrCh)
}

///
/// G1                      G2                          G3
/// DialContext()
///                         cc.dopts.balancer.Notify()
///                                                     cc.lbWatcher()
///                         <-doneChan
/// close()
/// -----------------------G2 leak------------------------------------
///

func main() {
	go DialContext() // G1
}



