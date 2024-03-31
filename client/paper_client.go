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

	AUTH uint8 = 2

	GET uint8 = 3
	SET uint8 = 4
	DEL uint8 = 5

	HAS uint8 = 6
	PEEK uint8 = 7
	TTL uint8 = 8
	SIZE uint8 = 9

	WIPE uint8 = 10

	RESIZE uint8 = 11
	POLICY uint8 = 12

	STATS uint8 = 13
)

const MAX_RECONNECT_ATTEMPTS = 3

type PaperClient struct {
	host string
	port uint32

	auth_token *string
	reconnect_attempts uint32

	tcp_client *tcp_client.TcpClient
}

func Connect(host string, port uint32) (*PaperClient, error) {
	tcp_client, err := tcp_client.Connect(host, port)

	if err != nil {
		return nil, err
	}

	var auth_token *string = nil
	var reconnect_attempts uint32 = 0

	client := PaperClient {
		host,
		port,

		auth_token,
		reconnect_attempts,

		tcp_client,
	}

	_, ping_err := client.Ping()

	if ping_err != nil {
		return nil, errors.New("Connection refused.")
	}

	return &client, nil
}

func (client *PaperClient) Disconnect() {
	client.tcp_client.GetConn().Close()
}

func (client *PaperClient) Ping() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PING)

	return client.process(sheet_writer)
}

func (client *PaperClient) Version() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(VERSION)

	return client.process(sheet_writer)
}

func (client *PaperClient) Auth(token string) (*response.Response[string], error) {
	client.auth_token = &token

	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(AUTH)
	sheet_writer.WriteString(token)

	return client.process(sheet_writer)
}

func (client *PaperClient) Get(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(GET)
	sheet_writer.WriteString(key)

	return client.process(sheet_writer)
}

func (client *PaperClient) Set(key string, value string, ttl uint32) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SET)
	sheet_writer.WriteString(key)
	sheet_writer.WriteString(value)
	sheet_writer.WriteU32(ttl)

	return client.process(sheet_writer)
}

func (client *PaperClient) Del(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(DEL)
	sheet_writer.WriteString(key)

	return client.process(sheet_writer)
}

func (client *PaperClient) Has(key string) (*response.Response[bool], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(HAS)
	sheet_writer.WriteString(key)

	return client.process_has(sheet_writer)
}

func (client *PaperClient) Peek(key string) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PEEK)
	sheet_writer.WriteString(key)

	return client.process(sheet_writer)
}

func (client *PaperClient) Ttl(key string, ttl uint32) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(TTL)
	sheet_writer.WriteString(key)
	sheet_writer.WriteU32(ttl)

	return client.process(sheet_writer)
}

func (client *PaperClient) Size(key string) (*response.Response[uint64], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SIZE)
	sheet_writer.WriteString(key)

	return client.process_size(sheet_writer)
}

func (client *PaperClient) Wipe() (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(WIPE)

	return client.process(sheet_writer)
}

func (client *PaperClient) Resize(size uint64) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(RESIZE)
	sheet_writer.WriteU64(size)

	return client.process(sheet_writer)
}

func (client *PaperClient) Policy(policy uint8) (*response.Response[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(POLICY)
	sheet_writer.WriteU8(policy)

	return client.process(sheet_writer)
}

func (client *PaperClient) Stats() (*response.Response[response.StatsData], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(STATS)

	return client.process_stats(sheet_writer)
}

func (client *PaperClient) reconnect() (error) {
	client.reconnect_attempts += 1

	if client.reconnect_attempts > MAX_RECONNECT_ATTEMPTS {
		return errors.New("Maximum reconnect attempts reached")
	}

	tcp_client, err := tcp_client.Connect(client.host, client.port)

	if err != nil {
		return err
	}

	client.tcp_client = tcp_client

	if client.auth_token != nil {
		if _, err := client.Auth(*client.auth_token); err != nil {
			return err
		}
	}

	return nil
}

func (client *PaperClient) process(sheet_writer *sheet_writer.SheetWriter) (*response.Response[string], error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process(sheet_writer)
	}

	response, err := client.get_str_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) process_has(sheet_writer *sheet_writer.SheetWriter) (*response.Response[bool], error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_has(sheet_writer)
	}

	response, err := client.get_has_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_has(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) process_size(sheet_writer *sheet_writer.SheetWriter) (*response.Response[uint64], error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_size(sheet_writer)
	}

	response, err := client.get_size_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_size(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) process_stats(sheet_writer *sheet_writer.SheetWriter) (*response.Response[response.StatsData], error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_stats(sheet_writer)
	}

	response, err := client.get_stats_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_stats(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) get_str_response() (*response.Response[string], error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	var ok_data *string
	var err_data *string

	ok_data = nil
	err_data = nil

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	if is_ok {
		ok_data = &data
	} else {
		err_data = &data
	}

	return response.New(ok_data, err_data), nil
}

func (client *PaperClient) get_has_response() (*response.Response[bool], error) {
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


func (client *PaperClient) get_size_response() (*response.Response[uint64], error) {
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

func (client *PaperClient) get_stats_response() (*response.Response[response.StatsData], error) {
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
