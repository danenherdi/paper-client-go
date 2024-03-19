package paper_client

import (
	"testing"
	"time"
	"math"
)

func TestPing(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Ping()

	if !response.IsOk() {
		t.Error("ping returned not ok")
	}

	if *response.Data() != "pong" {
		t.Error("ping did not return pong")
	}
}

func TestVersion(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Version()

	if !response.IsOk() {
		t.Error("version returned not ok")
	}

	if len(*response.Data()) == 0 {
		t.Error("version did not return version")
	}
}

func TestGetExistent(t *testing.T) {
	client := GetClient(t)

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
	client := GetClient(t)

	response, _ := client.Get("key")

	if response.IsOk() {
		t.Error("get returned ok for a key which does not exist")
	}

	if len(*response.ErrData()) == 0 {
		t.Error("get for key which does not exist did not return error message")
	}
}

func TestSetNoTtl(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Set("key", "value", 0)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with no ttl")
	}

	if *response.Data() != "done" {
		t.Error("set for a value with no ttl did not return \"done\"")
	}
}

func TestSetTtl(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Set("key", "value", 1)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with ttl")
	}

	if *response.Data() != "done" {
		t.Error("set for a value with ttl did not return \"done\"")
	}
}

func TestSetTtlExpiry(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Set("key", "value", 1)

	if !response.IsOk() {
		t.Error("set returned not ok for a value with ttl")
	}

	if *response.Data() != "done" {
		t.Error("set for a value with ttl did not return \"done\"")
	}

	got, _ := client.Get("key")

	if !got.IsOk() {
		t.Error("get returned not ok for a key which has not yet expired")
	}

	if *got.Data() != "value" {
		t.Errorf("get return %q instead of \"value\"", *response.Data())
	}

	time.Sleep(2 * time.Second)

	expired, _ := client.Get("key")

	if expired.IsOk() {
		t.Error("get returned ok for an expired key")
	}

	if len(*expired.ErrData()) == 0 {
		t.Error("get for expired key did not return an error message")
	}
}

func TestDelExistent(t *testing.T) {
	client := GetClient(t)

	client.Set("key", "value", 0)
	response, _ := client.Del("key")

	if !response.IsOk() {
		t.Error("del returned not ok for a key which exists")
	}

	if *response.Data() != "done" {
		t.Error("del for a key which exists did not return \"done\"")
	}
}

func TestDelNonExistent(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Del("key")

	if response.IsOk() {
		t.Error("del returned ok for a key which does not exist")
	}

	if len(*response.ErrData()) == 0 {
		t.Error("del for key which does not exist did not return error message")
	}
}

func TestHasExistent(t *testing.T) {
	client := GetClient(t)

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
	client := GetClient(t)

	response, _ := client.Has("key")

	if !response.IsOk() {
		t.Error("has returned not ok for a key which does not exist")
	}

	if *response.Data() {
		t.Error("has for key which does not exist returned true")
	}
}

func TestPeekExistent(t *testing.T) {
	client := GetClient(t)

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
	client := GetClient(t)

	response, _ := client.Peek("key")

	if response.IsOk() {
		t.Error("peek returned ok for a key which does not exist")
	}

	if len(*response.ErrData()) == 0 {
		t.Error("peek for key which does not exist did not return error message")
	}
}

func TestTtlExistent(t *testing.T) {
	client := GetClient(t)

	client.Set("key", "value", 0)
	response, _ := client.Ttl("key", 5)

	if !response.IsOk() {
		t.Error("ttl returned not ok for a key which exists")
	}

	if *response.Data() != "done" {
		t.Error("ttl for a key which exists did not return \"done\"")
	}
}

func TestTtlNonExistent(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Ttl("key", 5)

	if response.IsOk() {
		t.Error("ttl returned ok for a key which does not exist")
	}

	if len(*response.ErrData()) == 0 {
		t.Error("ttl for key which does not exist did not return error message")
	}
}

func TestSizeExistent(t *testing.T) {
	client := GetClient(t)

	client.Set("key", "value", 0)
	response, _ := client.Size("key")

	if !response.IsOk() {
		t.Error("size returned not ok for a key which exists")
	}

	if *response.Data() != 5 {
		t.Error("size for a key which exists did not return correct size")
	}
}

func TestSizeNonExistent(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Size("key")

	if response.IsOk() {
		t.Error("size returned ok for a key which does not exist")
	}

	if len(*response.ErrData()) == 0 {
		t.Error("size for key which does not exist did not return error message")
	}
}

func TestWipe(t *testing.T) {
	client := GetClient(t)

	client.Set("key", "value", 0)
	response, _ := client.Wipe()

	if !response.IsOk() {
		t.Error("wipe returned not ok")
	}

	if *response.Data() != "done" {
		t.Error("wipe did not return \"done\"")
	}

	got, _ := client.Get("key")

	if got.IsOk() {
		t.Error("get returned ok for wiped key")
	}

	if len(*got.ErrData()) == 0 {
		t.Error("get for wiped key did not return error message")
	}
}

func TestResize(t *testing.T) {
	var INITIAL_SIZE = uint64(10 * math.Pow(1024, 2))
	var UPDATED_SIZE = uint64(20 * math.Pow(1024, 2))

	client := GetClient(t)

	initial, _ := client.Resize(INITIAL_SIZE)

	if !initial.IsOk() {
		t.Error("resize returned not ok")
	}

	if *initial.Data() != "done" {
		t.Error("resize did not return \"done\"")
	}

	if GetCacheSize() != INITIAL_SIZE {
		t.Error("cache has incorrect initial size")
	}

	updated, _ := client.Resize(UPDATED_SIZE)

	if !updated.IsOk() {
		t.Error("resize returned not ok")
	}

	if *updated.Data() != "done" {
		t.Error("resize did not return \"done\"")
	}

	if GetCacheSize() != UPDATED_SIZE {
		t.Error("cache has incorrect updated size")
	}
}

func TestPolicy(t *testing.T) {
	var INITIAL_POLICY = POLICY_LFU
	var UPDATED_POLICY = POLICY_LRU

	client := GetClient(t)

	initial, _ := client.Policy(INITIAL_POLICY)

	if !initial.IsOk() {
		t.Error("policy returned not ok")
	}

	if *initial.Data() != "done" {
		t.Error("policy did not return \"done\"")
	}

	if GetCachePolicy() != INITIAL_POLICY {
		t.Error("cache has incorrect initial policy")
	}

	updated, _ := client.Policy(UPDATED_POLICY)

	if !updated.IsOk() {
		t.Error("policy returned not ok")
	}

	if *updated.Data() != "done" {
		t.Error("policy did not return \"done\"")
	}

	if GetCachePolicy() != UPDATED_POLICY {
		t.Error("cache has incorrect updated policy")
	}
}

func TestStats(t *testing.T) {
	client := GetClient(t)

	response, _ := client.Stats()

	if !response.IsOk() {
		t.Error("stats returned not ok")
	}
}

var client, _ = Connect("127.0.0.1", 3145)

func GetClient(t *testing.T) (*PaperClient) {
	client.Wipe()
	return client
}

func GetCacheSize() (uint64) {
	response, _ := client.Stats()
	return (*response.Data()).MaxSize();
}

func GetCachePolicy() (uint8) {
	response, _ := client.Stats()
	return (*response.Data()).Policy();
}
