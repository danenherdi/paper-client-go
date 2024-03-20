package paper_pool

import (
	"testing"
)

func TestClient(t *testing.T) {
	pool, _ := Connect("127.0.0.1", 3145, 2)

	for i := 0; i < 10; i++ {
		lockable_client := pool.LockableClient()

		client := lockable_client.Lock()
		response, _ := client.Ping()

		if !response.IsOk() {
			t.Error("pool client ping returned not ok")
		}

		if *response.Data() != "pong" {
			t.Error("pool client ping did not return pong")
		}

		lockable_client.Unlock()
	}
}
