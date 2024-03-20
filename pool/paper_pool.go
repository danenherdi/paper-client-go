package paper_pool

import (
	"sync"
	"sync/atomic"
	"github.com/Paper-Book/paper-client-go/client"
)

type PaperPool struct {
	clients []*LockableClient
	index uint32
}

type LockableClient struct {
	client *paper_client.PaperClient
	lock *sync.Mutex
}

func Connect(host string, port uint32, size uint32) (*PaperPool, error) {
	clients := []*LockableClient{}

	for i := uint32(0); i<size; i++ {
		client, err := paper_client.Connect(host, port)

		if err != nil {
			return nil, err
		}

		lock := &sync.Mutex{}

		locked_client := LockableClient {
			client,
			lock,
		}

		clients = append(clients, &locked_client)
	}

	index := uint32(0)

	pool := PaperPool {
		clients,
		index,
	}

	return &pool, nil
}

func (pool *PaperPool) LockableClient() (*LockableClient) {
	client := pool.clients[pool.index]

	new_index := (pool.index + 1) % uint32(len(pool.clients))
	atomic.StoreUint32(&pool.index, new_index)

	return client
}

func (lockable_client *LockableClient) Lock() (*paper_client.PaperClient) {
	lockable_client.lock.Lock()
	return lockable_client.client
}

func (lockable_client *LockableClient) Unlock() {
	lockable_client.lock.Unlock()
}
