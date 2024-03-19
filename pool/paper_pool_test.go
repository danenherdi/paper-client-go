package paper_pool

import (
	"testing"
)

func TestClient(t *testing.T) {
	pool, _ := Connect("127.0.0.1", 3145, 2)
	client := pool.Client()

	response, _ := client.Ping()

	if !response.IsOk() {
		t.Error("pool client ping returned not ok")
	}

	if *response.Data() != "pong" {
		t.Error("pool client ping did not return pong")
	}
}
