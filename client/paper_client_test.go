/*
 * Copyright (c) Kia Shakiba
 *
 * This source code is licensed under the GNU AGPLv3 license found in the
 * LICENSE file in the root directory of this source tree.
 */

package paper_client

import (
	"testing"
	"time"
	"math"
)

func TestPing(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Ping()

	if !response.IsOk() {
		t.Error("ping returned not ok")
	}

	if *response.Data() != "pong" {
		t.Error("ping did not return pong")
	}
}

func TestVersion(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Version()

	if !response.IsOk() {
		t.Error("version returned not ok")
	}

	if len(*response.Data()) == 0 {
		t.Error("version did not return version")
	}
}

func TestAuthIncorrect(t *testing.T) {
	client := InitClient(t, false)
	defer client.Disconnect()

	response, _ := client.Auth("incorrect_auth_token")

	if response.IsOk() {
		t.Error("auth with incorrect token returned ok")
	}
}

func TestAuthCorrect(t *testing.T) {
	client := InitClient(t, false)
	defer client.Disconnect()

	response, _ := client.Auth("auth_token")

	if !response.IsOk() {
		t.Error("auth with correct token returned not ok")
	}
}

func TestGetExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Get("key")

	if !response.IsOk() {
		t.Error("get returned not ok for a key which exists")
	}

	if *response.Data() != "value" {
		t.Errorf("get return %q instead of \"value\"", *response.Data())
	}
}

func TestGetNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Get("key")

	if response.IsOk() {
		t.Error("get returned ok for a key which does not exist")
	}

	if response.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("get for key which does not exist did not return correct error message")
	}
}

func TestSetNoTtl(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Set("key", "value", 0)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with no ttl")
	}
}

func TestSetTtl(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Set("key", "value", 1)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with ttl")
	}
}

func TestSetTtlExpiry(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Set("key", "value", 1)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with ttl")
	}

	got, _ := client.Get("key")

	if !got.IsOk() {
		t.Error("get returned not ok for a key which has not yet expired")
	}

	if *got.Data() != "value" {
		t.Errorf("get return %q instead of \"value\"", *got.Data())
	}

	time.Sleep(2 * time.Second)

	expired, _ := client.Get("key")

	if expired.IsOk() {
		t.Error("get returned ok for an expired key")
	}

	if expired.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("get for expired key did not return a correct error")
	}
}

func TestDelExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Del("key")

	if !response.IsOk() {
		t.Error("del returned not ok for a key which exists")
	}
}

func TestDelNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Del("key")

	if response.IsOk() {
		t.Error("del returned ok for a key which does not exist")
	}

	if response.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("del for key which does not exist did not return correct error")
	}
}

func TestHasExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Has("key")

	if !response.IsOk() {
		t.Error("has returned not ok for a key which exists")
	}

	if !*response.Data() {
		t.Error("has for a key which exists returned false")
	}
}

func TestHasNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Has("key")

	if !response.IsOk() {
		t.Error("has returned not ok for a key which does not exist")
	}

	if *response.Data() {
		t.Error("has for key which does not exist returned true")
	}
}

func TestPeekExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Peek("key")

	if !response.IsOk() {
		t.Error("peek returned not ok for a key which exists")
	}

	if *response.Data() != "value" {
		t.Errorf("peek return %q instead of \"value\"", *response.Data())
	}
}

func TestPeekNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Peek("key")

	if response.IsOk() {
		t.Error("peek returned ok for a key which does not exist")
	}

	if response.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("peek for key which does not exist did not return correct error")
	}
}

func TestTtlExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Ttl("key", 5)

	if !response.IsOk() {
		t.Error("ttl returned not ok for a key which exists")
	}
}

func TestTtlNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Ttl("key", 5)

	if response.IsOk() {
		t.Error("ttl returned ok for a key which does not exist")
	}

	if response.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("ttl for key which does not exist did not return correct error")
	}
}

func TestSizeExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Size("key")

	if !response.IsOk() {
		t.Error("size returned not ok for a key which exists")
	}

	if *response.Data() == 0 {
		t.Error("size for a key which exists did not return correct size")
	}
}

func TestSizeNonExistent(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Size("key")

	if response.IsOk() {
		t.Error("size returned ok for a key which does not exist")
	}

	if response.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("size for key which does not exist did not return correct error")
	}
}

func TestWipe(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, _ := client.Wipe()

	if !response.IsOk() {
		t.Error("wipe returned not ok")
	}

	got, _ := client.Get("key")

	if got.IsOk() {
		t.Error("get returned ok for wiped key")
	}

	if got.Error() != PAPER_ERROR_KEY_NOT_FOUND {
		t.Error("get for wiped key did not return correct error")
	}
}

func TestResize(t *testing.T) {
	var INITIAL_SIZE = uint64(10 * math.Pow(1024, 2))
	var UPDATED_SIZE = uint64(20 * math.Pow(1024, 2))

	client := InitClient(t, true)
	defer client.Disconnect()

	initial, _ := client.Resize(INITIAL_SIZE)

	if !initial.IsOk() {
		t.Error("resize returned not ok")
	}

	if GetCacheSize(client) != INITIAL_SIZE {
		t.Error("cache has incorrect initial size")
	}

	updated, _ := client.Resize(UPDATED_SIZE)

	if !updated.IsOk() {
		t.Error("resize returned not ok")
	}

	if GetCacheSize(client) != UPDATED_SIZE {
		t.Error("cache has incorrect updated size")
	}
}

func TestPolicy(t *testing.T) {
	var INITIAL_POLICY = "lfu"
	var UPDATED_POLICY = "lru"

	client := InitClient(t, true)
	defer client.Disconnect()

	initial, _ := client.Policy(INITIAL_POLICY)

	if !initial.IsOk() {
		t.Error("policy returned not ok")
	}

	if GetCachePolicy(client) != INITIAL_POLICY {
		t.Error("cache has incorrect initial policy")
	}

	updated, _ := client.Policy(UPDATED_POLICY)

	if !updated.IsOk() {
		t.Error("policy returned not ok")
	}

	if GetCachePolicy(client) != UPDATED_POLICY {
		t.Error("cache has incorrect updated policy")
	}
}

func TestStatus(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	response, _ := client.Status()

	if !response.IsOk() {
		t.Error("status returned not ok")
	}
}

func TestReconnect(t *testing.T) {
	client := InitClient(t, true)
	defer client.Disconnect()

	pre_disconnect, _ := client.Has("key")

	if !pre_disconnect.IsOk() {
		t.Error("has returned not ok")
	}

	client.Disconnect()

	post_disconnect, err := client.Has("key")

	if err != nil {
		t.Error("an error occurred after disconnect")
	}

	if !post_disconnect.IsOk() {
		t.Error("has returned not ok")
	}
}

func InitClient(t *testing.T, authorize bool) (*PaperClient) {
	client, err := Connect("paper://127.0.0.1:3145")

	if err != nil {
		t.Error("Could not connect client")
	}

	if authorize {
		response, _ := client.Auth("auth_token")

		if !response.IsOk() {
			t.Error("Could not authorize client")
		}
	}

	client.Wipe()

	return client
}

func GetCacheSize(client *PaperClient) (uint64) {
	response, _ := client.Status()
	return (*response.Data()).MaxSize();
}

func GetCachePolicy(client *PaperClient) (string) {
	response, _ := client.Status()
	return (*response.Data()).Policy();
}
