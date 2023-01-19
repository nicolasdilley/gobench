/*
 * Project: moby
 * Issue or PR  : https://github.com/moby/moby/pull/33293
 * Buggy version: 4921171587c09d0fcd8086a62a25813332f44112
 * fix commit-id:
 * Flaky: 100/100
 */
package main

import (
	"errors"
	"math/rand"

)

func MayReturnError() error {
	if rand.Int31n(2) >= 1 {
		return errors.New("Error")
	}
	return nil
}
func containerWait() error {
	errC := make(chan error)
	err := MayReturnError()
	if err != nil {
		errC <- err /// Block here
		return err
	}
	return err
}

///
/// G1
/// containerWait()
/// errC <- err
/// ---------G1 leak---------------
///

func main() {
	go func() { // G1
		err := containerWait()
		if err != nil {
			return
		}
	}()
}



