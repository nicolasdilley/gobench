package cockroach25456

import (
	"testing"
)

func process(quiescer chan struct{}) {
	<-quiescer
}

func TestCockroach25456(t *testing.T) {
	quiescer := make(chan struct{})

	for i := 0; i < 2; i++ {
		process(quiescer)
	}
}
