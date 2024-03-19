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
	TTL uint8 = 7
	SIZE uint8 = 8

	WIPE uint8 = 9

	RESIZE uint8 = 10
	POLICY uint8 = 11

	STATS uint8 = 12
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *bool
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadBool()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
}

func (client *PaperClient) Ttl(key string, ttl uint32) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(TTL)
	sheet_writer.WriteString(key)
	sheet_writer.WriteU32(ttl)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
}

func (client *PaperClient) Size(key string) (*response.Response[uint64], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SIZE)
	sheet_writer.WriteString(key)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	var ok_data *uint64
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
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

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		ok_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
}

func (client *PaperClient) Stats() (*response.Response[response.StatsData], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(STATS)

	client.tcp_client.Send(sheet_writer)

	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	var ok_data *response.StatsData
	var err_data *string

	ok_data = nil
	err_data = nil

	if is_ok {
		max_size, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		used_size, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		total_gets, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		total_sets, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		total_dels, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		miss_ratio, err := sheet_reader.ReadF64()

		if err != nil {
			return nil, err
		}

		policy, err := sheet_reader.ReadU8()

		if err != nil {
			return nil, err
		}

		uptime, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		ok_response_data := response.NewStatsData(
			max_size,
			used_size,

			total_gets,
			total_sets,
			total_dels,

			miss_ratio,

			policy,
			uptime,
		)

		ok_data = &ok_response_data
	} else {
		err_response_data, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		err_data = &err_response_data
	}

	return response.New(ok_data, err_data), nil
}
