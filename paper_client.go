package paper_client

import (
	"errors"
	"internal/sheet_writer"
	"internal/sheet_reader"
	"internal/tcp_client"
	"internal/response"
)

const (
	POLICY_LFU uint8 = 0
	POLICY_FIFO uint8 = 1
	POLICY_LRU uint8 = 2
	POLICY_MRU uint8 = 3
)

const (
	PING uint8 = 0
	VERSION uint8 = 1

	GET uint8 = 2
	SET uint8 = 3
	DEL uint8 = 4

	HAS uint8 = 5
	PEEK uint8 = 6

	WIPE uint8 = 7

	RESIZE uint8 = 8
	POLICY uint8 = 9

	STATS uint8 = 10
)

type PaperClient struct {
	tcp_client *tcp_client.TcpClient
}

func Connect(host string, port uint32) (*PaperClient, error) {
	tcp_client, err := tcp_client.Connect(host, port)

	if err != nil {
		return nil, err
	}

	client := PaperClient {
		tcp_client,
	}

	_, ping_err := client.Ping()

	if ping_err != nil {
		return nil, errors.New("Connection refused.")
	}

	return &client, nil
}

func (client *PaperClient) Ping() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PING)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Version() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(VERSION)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Get(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(GET)
	sheet_writer.WriteString(key)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Set(key string, value string, ttl uint32) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SET)
	sheet_writer.WriteString(key)
	sheet_writer.WriteString(value)
	sheet_writer.WriteU32(ttl)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Del(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(DEL)
	sheet_writer.WriteString(key)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Has(key string) (*response.Response[bool], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(HAS)
	sheet_writer.WriteString(key)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Peek(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PEEK)
	sheet_writer.WriteString(key)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Wipe() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(WIPE)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Resize(size uint64) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(RESIZE)
	sheet_writer.WriteU64(size)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Policy(policy uint8) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(POLICY)
	sheet_writer.WriteU8(policy)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.New(is_ok, data), nil
}

func (client *PaperClient) Stats() (*response.Response[response.StatsData], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(STATS)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	_, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	max_size, _ := sheet_reader.ReadU64()
	used_size, _ := sheet_reader.ReadU64()

	total_gets, _ := sheet_reader.ReadU64()
	total_sets, _ := sheet_reader.ReadU64()
	total_dels, _ := sheet_reader.ReadU64()

	miss_ratio, _ := sheet_reader.ReadF64()

	policy, _ := sheet_reader.ReadU8()
	uptime, _ := sheet_reader.ReadU64()

	stats_response := response.NewStatsResponse(
		max_size,
		used_size,

		total_gets,
		total_sets,
		total_dels,

		miss_ratio,

		policy,
		uptime,
	)

	return stats_response, nil
}
