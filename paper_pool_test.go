/*
 * Copyright (c) Kia Shakiba
 *
 * This source code is licensed under the GNU AGPLv3 license found in the
 * LICENSE file in the root directory of this source tree.
 */

package paperclient

import (
	"testing"
)

func TestClient(t *testing.T) {
	pool, _ := PoolConnect("paper://127.0.0.1:3145", 2)
	defer pool.Disconnect()

	for i := 0; i < 10; i++ {
		lockable_client := pool.LockableClient()

		client := lockable_client.Lock()
		response, err := client.Ping()

		if err != nil {
			t.Error("pool client ping returned an error")
		}

		if response != "pong" {
			t.Error("pool client ping did not return pong")
		}

		lockable_client.Unlock()
	}
}

func TestAuthInvalid(t *testing.T) {
	pool, _ := PoolConnect("paper://127.0.0.1:3145", 2)
	defer pool.Disconnect()

	lockable_client := pool.LockableClient()

	client := lockable_client.Lock()
	err := client.Set("key", "value", 0)

	if err == nil {
		t.Error("unauthorized pool client did not return an error")
	}

	lockable_client.Unlock()
}

func TestAuthValid(t *testing.T) {
	pool, _ := PoolConnect("paper://127.0.0.1:3145", 2)
	defer pool.Disconnect()

	pool.Auth("auth_token")

	lockable_client := pool.LockableClient()

	client := lockable_client.Lock()
	err := client.Set("key", "value", 0)

	if err != nil {
		t.Error("authorized pool client returned an error")
	}

	lockable_client.Unlock()
}
