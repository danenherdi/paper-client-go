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
	client := initClient(t, true)
	defer client.Disconnect()

	response, err := client.Ping()

	if err != nil {
		t.Error(err)
	}

	if response != "pong" {
		t.Error("ping did not return pong")
	}
}

func TestVersion(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	response, err := client.Version()

	if err != nil {
		t.Error(err)
	}

	if len(response) == 0 {
		t.Error("version did not return version")
	}
}

func TestAuthIncorrect(t *testing.T) {
	client := initClient(t, false)
	defer client.Disconnect()

	err := client.Auth("incorrect_auth_token")

	if err == nil {
		t.Error("incorrect auth token did not return an error")
	}
}

func TestAuthCorrect(t *testing.T) {
	client := initClient(t, false)
	defer client.Disconnect()

	err := client.Auth("auth_token")

	if err != nil {
		t.Error("correct auth token returned an error")
	}
}

func TestGetExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, err := client.Get("key")

	if err != nil {
		t.Error("get returned an error for a key which exists")
	}

	if response != "value" {
		t.Errorf("get return %q instead of \"value\"", response)
	}
}

func TestGetNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	_, err := client.Get("key")

	if err == nil {
		t.Error("get did not return an error for a key which does not exist")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("get for key which does not exist did not return correct error")
	}
}

func TestSetNoTtl(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Set("key", "value", 0)

	if err != nil {
		t.Error("set returned an error for a value with no ttl")
	}
}

func TestSetTtl(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Set("key", "value", 1)

	if err != nil {
		t.Error("set returned an error for a value with ttl")
	}
}

func TestSetTtlExpiry(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Set("key", "value", 1)

	if err != nil {
		t.Error("set returned an error for a value with ttl")
	}

	got, err := client.Get("key")

	if err != nil {
		t.Error("get returned an error for a key which has not yet expired")
	}

	if got != "value" {
		t.Errorf("get return %q instead of \"value\"", got)
	}

	time.Sleep(2 * time.Second)

	_, err = client.Get("key")

	if err == nil {
		t.Error("get did not return an error for an expired key")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("get for expired key did not return a correct error")
	}
}

func TestDelExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	err := client.Del("key")

	if err != nil {
		t.Error("del returned an error for a key which exists")
	}
}

func TestDelNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Del("key")

	if err == nil {
		t.Error("del did not return an error for a key which does not exist")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("del for key which does not exist did not return correct error")
	}
}

func TestHasExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	has, err := client.Has("key")

	if err != nil {
		t.Error("has returned an error for a key which exists")
	}

	if !has {
		t.Error("has for a key which exists returned false")
	}
}

func TestHasNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	has, err := client.Has("key")

	if err != nil {
		t.Error("has returned an error for a key which does not exist")
	}

	if has {
		t.Error("has for key which does not exist returned true")
	}
}

func TestPeekExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	response, err := client.Peek("key")

	if err != nil {
		t.Error("peek returned an error for a key which exists")
	}

	if response != "value" {
		t.Errorf("peek return %q instead of \"value\"", response)
	}
}

func TestPeekNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	_, err := client.Peek("key")

	if err == nil {
		t.Error("peek did not return an error for a key which does not exist")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("peek for key which does not exist did not return correct error")
	}
}

func TestTtlExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	err := client.Ttl("key", 5)

	if err != nil {
		t.Error("ttl returned an error for a key which exists")
	}
}

func TestTtlNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Ttl("key", 5)

	if err == nil {
		t.Error("ttl did not return an error for a key which does not exist")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("ttl for key which does not exist did not return correct error")
	}
}

func TestSizeExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	size, err := client.Size("key")

	if err != nil {
		t.Error("size returned an error for a key which exists")
	}

	if size == 0 {
		t.Error("size for a key which exists did not return correct size")
	}
}

func TestSizeNonExistent(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	_, err := client.Size("key")

	if err == nil {
		t.Error("size did not return an error for a key which does not exist")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("size for key which does not exist did not return correct error")
	}
}

func TestWipe(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	client.Set("key", "value", 0)
	err := client.Wipe()

	if err != nil {
		t.Error("wipe returned an error")
	}

	_, err = client.Get("key")

	if err == nil {
		t.Error("get dit not return an error for wiped key")
	}

	if err != PaperErrorKeyNotFound {
		t.Error("get for wiped key did not return correct error")
	}
}

func TestResize(t *testing.T) {
	var INITIAL_SIZE = uint64(10 * math.Pow(1024, 2))
	var UPDATED_SIZE = uint64(20 * math.Pow(1024, 2))

	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Resize(INITIAL_SIZE)

	if err != nil {
		t.Error("resize returned an error")
	}

	if getCacheSize(client) != INITIAL_SIZE {
		t.Error("cache has incorrect initial size")
	}

	err = client.Resize(UPDATED_SIZE)

	if err != nil {
		t.Error("resize returned an error")
	}

	if getCacheSize(client) != UPDATED_SIZE {
		t.Error("cache has incorrect updated size")
	}
}

func TestPolicy(t *testing.T) {
	var INITIAL_POLICY = "lfu"
	var UPDATED_POLICY = "lru"

	client := initClient(t, true)
	defer client.Disconnect()

	err := client.Policy(INITIAL_POLICY)

	if err != nil {
		t.Error("policy returned an error")
	}

	if getCachePolicy(client) != INITIAL_POLICY {
		t.Error("cache has incorrect initial policy")
	}

	err = client.Policy(UPDATED_POLICY)

	if err != nil {
		t.Error("policy returned an error")
	}

	if getCachePolicy(client) != UPDATED_POLICY {
		t.Error("cache has incorrect updated policy")
	}
}

func TestStatus(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	_, err := client.Status()

	if err != nil {
		t.Error("status returned an error")
	}
}

func TestReconnect(t *testing.T) {
	client := initClient(t, true)
	defer client.Disconnect()

	_, err := client.Has("key")

	if err != nil {
		t.Error("has returned an error before disconnect")
	}

	client.Disconnect()

	_, err = client.Has("key")

	if err != nil {
		t.Error("has returned an error after disconnect")
	}
}

func initClient(t *testing.T, authorize bool) (*PaperClient) {
	client, err := ClientConnect("paper://127.0.0.1:3145")

	if err != nil {
		t.Error("Could not connect client")
	}

	if authorize {
		err := client.Auth("auth_token")

		if err != nil {
			t.Error("Could not authorize client")
		}

		client.Wipe()
	}

	return client
}

func getCacheSize(client *PaperClient) (uint64) {
	status, _ := client.Status()
	return status.max_size
}

func getCachePolicy(client *PaperClient) (string) {
	status, _ := client.Status()
	return status.policy
}
