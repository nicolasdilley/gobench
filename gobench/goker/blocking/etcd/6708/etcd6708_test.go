package etcd6708

import (
	"context"
	"sync"
	"testing"
)

type EndpointSelectionMode int

const (
	EndpointSelectionRandom EndpointSelectionMode = iota
	EndpointSelectionPrioritizeLeader
)

type Client interface {
	Sync(ctx context.Context)
	SetEndpoints()
	httpClient
}

type httpClient interface {
	Do(context.Context)
}

type httpClusterClient struct {
	mu            sync.RWMutex
	selectionMode EndpointSelectionMode
}

func (c *httpClusterClient) getLeaderEndpoint() {
	c.Do(context.Background())
}

func (c *httpClusterClient) SetEndpoints() {
	switch c.selectionMode {
	case EndpointSelectionRandom:
	case EndpointSelectionPrioritizeLeader:
		c.getLeaderEndpoint()
	}
}

func (c *httpClusterClient) Do(ctx context.Context) {
	c.mu.RLock() // block here
	c.mu.RUnlock()
}

func (c *httpClusterClient) Sync(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.SetEndpoints()
}

type httpMembersAPI struct {
	client httpClient
}

func TestEtcd6708(t *testing.T) {
	hc := &httpClusterClient{
		selectionMode: EndpointSelectionPrioritizeLeader,
	}
	hc.Sync(context.Background())
}
