/*
 * Project: kubernetes
 * Issue or PR  : https://github.com/kubernetes/kubernetes/pull/6632
 * Buggy version: e597b41d939573502c8dda1dde7bf3439325fb5d
 * fix commit-id: 82afb7ab1fe12cf2efceede2322d082eaf5d5adc
 * Flaky: 4/100
 * Description:
 *   This is a lock-channel bug. When resetChan is full, WriteFrame
 * holds the lock and blocks on the channel. Then monitor() fails
 * to close the resetChan because lock is already held by WriteFrame.
 *   Fix: create a goroutine to drain the channel
 */
package main

import (
	"sync"

)

type Connection struct {
	closeChan chan bool
}

type idleAwareFramer struct {
	resetChan chan bool
	writeLock sync.Mutex
	conn      *Connection
}

func (i *idleAwareFramer) monitor() {
Loop:
	for {
		select {
		case <-i.conn.closeChan:
			i.writeLock.Lock()
			close(i.resetChan)
			i.resetChan = nil
			i.writeLock.Unlock()
			break Loop
		}
	}
}

func (i *idleAwareFramer) WriteFrame() {
	i.writeLock.Lock()
	defer i.writeLock.Unlock()
	if i.resetChan == nil {
		return
	}
	i.resetChan <- true
}

///
/// G1						G2					helper goroutine
/// i.monitor()
/// <-i.conn.closeChan
///							i.WriteFrame()
///							i.writeLock.Lock()
///							i.resetChan <-
///												i.conn.closeChan<-
///	i.writeLock.Lock()
///	----------------------G1,G2 deadlock------------------------
///
func main() {
	i := &idleAwareFramer{
		resetChan: make(chan bool),
		conn: &Connection{
			closeChan: make(chan bool),
		},
	}

	go func() { // helper goroutine
		i.conn.closeChan <- true
	}()
	go i.monitor()    // G1
	go i.WriteFrame() // G2
}



