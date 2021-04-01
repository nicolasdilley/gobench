package cockroach25456

import (
	"testing"
)

type Stopper struct {
	quiescer chan struct{}
}

type Store struct {
	stopper          *Stopper
	consistencyQueue *consistencyQueue
}

type Replica struct {
	store *Store
}

type consistencyQueue struct{}

func (q *consistencyQueue) process(repl *Replica) {
	<-repl.store.stopper.quiescer
}

type testContext struct {
	store *Store
	repl  *Replica
}

func (tc *testContext) StartWithStoreConfig(stopper *Stopper) {
	if tc.store == nil {
		tc.store = &Store{
			consistencyQueue: &consistencyQueue{},
		}
	}
	tc.store.stopper = stopper
	tc.repl = &Replica{store: tc.store}
}

func TestCockroach25456(t *testing.T) {
	stopper := &Stopper{quiescer: make(chan struct{})}
	tc := testContext{}
	tc.StartWithStoreConfig(stopper)

	for i := 0; i < 2; i++ {
		tc.store.consistencyQueue.process(tc.repl)
	}
}
