package paper_pool

import (
	"sync"
	"sync/atomic"
	"github.com/Paper-Book/paper-client-go/client"
)

type PaperPool struct {
	clients []*paper_client.PaperClient
	locks []*sync.Mutex

	index uint32
}

func Connect(host string, port uint32, size uint32) (*PaperPool, error) {
	clients := []*paper_client.PaperClient{}
	locks := []*sync.Mutex{}

	for i := uint32(0); i<size; i++ {
		client, err := paper_client.Connect(host, port)

		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
		locks = append(locks, &sync.Mutex{})
	}

	index := uint32(0)

	pool := PaperPool {
		clients,
		locks,
		index,
	}

	return &pool, nil
}

func (pool *PaperPool) Client() (*paper_client.PaperClient) {
	client := pool.clients[pool.index]

	new_index := (pool.index + 1) % uint32(len(pool.clients))
	atomic.StoreUint32(&pool.index, new_index)

	return client
}
