package paper_client

import (
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

	return &client, nil
}

func (client *PaperClient) Ping() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(0)

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
	sheet_writer.WriteU8(1)

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
	sheet_writer.WriteU8(2)
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
	sheet_writer.WriteU8(3)
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
	sheet_writer.WriteU8(4)
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
	sheet_writer.WriteU8(5)

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
	sheet_writer.WriteU8(6)
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
	sheet_writer.WriteU8(7)
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
	sheet_writer.WriteU8(8)

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
