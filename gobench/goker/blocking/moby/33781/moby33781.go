/*
 * Project: moby
 * Issue or PR  : https://github.com/moby/moby/pull/33781
 * Buggy version: 33fd3817b0f5ca4b87f0a75c2bd583b4425d392b
 * fix commit-id: 67297ba0051d39be544009ba76abea14bc0be8a4
 * Flaky: 25/100
 * Description:
 *   The goroutine created using anonymous function is blocked at
 * sending message to a unbuffered channel. However there exists a
 * path in the parent goroutine where the parent function will
 * return without draining the channel.
 */

package main

import (
	"context"

	"time"
)

func monitor(stop chan bool) {
	probeInterval := 50 * time.Nanosecond
	probeTimeout := 50 * time.Nanosecond
	results := make(chan bool)
	for {
		select {
		case <-stop:
			return
		case <-time.After(probeInterval):
			ctx, cancelProbe := context.WithTimeout(context.Background(), probeTimeout)
			go func() { // G3
				results <- true
				close(results)
			}()
			select {
			case <-stop:
				// results should be drained here
				cancelProbe()
				return
			case <-results:
				cancelProbe()
			case <-ctx.Done():
				cancelProbe()
				<-results
			}
		}
	}
}

///
/// G1 				G2				G3
/// monitor()
/// <-time.After()
/// 				stop <-
/// <-stop
/// 				return
/// cancelProbe()
/// return
/// 								result<-
///----------------G3 leak------------------
///

func main() {
	stop := make(chan bool)
	go monitor(stop) // G1
	go func() {      // G2
		time.Sleep(50 * time.Nanosecond)
		stop <- true
	}()
}



