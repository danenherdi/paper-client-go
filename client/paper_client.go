package paper_client

import (
	"errors"
	"strings"
	"internal/sheet_writer"
	"internal/sheet_reader"
	"internal/tcp_client"
	"internal/response"
)

const (
	PAPER_ERROR_INTERNAL uint8 					= 0

	PAPER_ERROR_UNREACHABLE_SERVER uint8		= 1
	PAPER_ERROR_MAX_CONNECTIONS_EXCEEDED uint8	= 2
	PAPER_ERROR_UNAUTHORIZED uint8				= 3

	PAPER_ERROR_KEY_NOT_FOUND uint8				= 4

	PAPER_ERROR_ZERO_VALUE_SIZE uint8			= 5
	PAPER_ERROR_EXCEEDING_VALUE_SIZE uint8		= 6

	PAPER_ERROR_UNCONFIGURED_POLICY uint8		= 7
	PAPER_ERROR_INVALID_POLICY uint8			= 8

	PAPER_ERROR_ZERO_CACHE_SIZE uint8			= 9
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
	addr string

	auth_token *string
	reconnect_attempts uint32

	tcp_client *tcp_client.TcpClient
}

func Connect(paper_addr string) (*PaperClient, error) {
	addr_ptr, err := parse_paper_addr(paper_addr)

	if err != nil {
		return nil, err
	}

	addr := *addr_ptr
	tcp_client, err := tcp_client.Connect(addr)

	if err != nil {
		return nil, err
	}

	var auth_token *string = nil
	var reconnect_attempts uint32 = 0

	client := PaperClient {
		addr,

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

func (client *PaperClient) Ping() (*response.DataResponse[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PING)

	return client.process_data(sheet_writer)
}

func (client *PaperClient) Version() (*response.DataResponse[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(VERSION)

	return client.process_data(sheet_writer)
}

func (client *PaperClient) Auth(token string) (*response.Response, error) {
	client.auth_token = &token

	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(AUTH)
	sheet_writer.WriteString(token)

	return client.process(sheet_writer)
}

func (client *PaperClient) Get(key string) (*response.DataResponse[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(GET)
	sheet_writer.WriteString(key)

	return client.process_data(sheet_writer)
}

func (client *PaperClient) Set(key string, value string, ttl uint32) (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SET)
	sheet_writer.WriteString(key)
	sheet_writer.WriteString(value)
	sheet_writer.WriteU32(ttl)

	return client.process(sheet_writer)
}

func (client *PaperClient) Del(key string) (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(DEL)
	sheet_writer.WriteString(key)

	return client.process(sheet_writer)
}

func (client *PaperClient) Has(key string) (*response.DataResponse[bool], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(HAS)
	sheet_writer.WriteString(key)

	return client.process_has(sheet_writer)
}

func (client *PaperClient) Peek(key string) (*response.DataResponse[string], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(PEEK)
	sheet_writer.WriteString(key)

	return client.process_data(sheet_writer)
}

func (client *PaperClient) Ttl(key string, ttl uint32) (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(TTL)
	sheet_writer.WriteString(key)
	sheet_writer.WriteU32(ttl)

	return client.process(sheet_writer)
}

func (client *PaperClient) Size(key string) (*response.DataResponse[uint32], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(SIZE)
	sheet_writer.WriteString(key)

	return client.process_size(sheet_writer)
}

func (client *PaperClient) Wipe() (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(WIPE)

	return client.process(sheet_writer)
}

func (client *PaperClient) Resize(size uint64) (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(RESIZE)
	sheet_writer.WriteU64(size)

	return client.process(sheet_writer)
}

func (client *PaperClient) Policy(policy string) (*response.Response, error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(POLICY)
	sheet_writer.WriteString(policy)

	return client.process(sheet_writer)
}

func (client *PaperClient) Stats() (*response.DataResponse[response.StatsData], error) {
	sheet_writer := sheet_writer.New()
	sheet_writer.WriteU8(STATS)

	return client.process_stats(sheet_writer)
}

func (client *PaperClient) reconnect() (error) {
	client.reconnect_attempts += 1

	if client.reconnect_attempts > MAX_RECONNECT_ATTEMPTS {
		return errors.New("Maximum reconnect attempts reached")
	}

	tcp_client, err := tcp_client.Connect(client.addr)

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

func (client *PaperClient) process(sheet_writer *sheet_writer.SheetWriter) (*response.Response, error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process(sheet_writer)
	}

	response, err := client.get_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) process_data(sheet_writer *sheet_writer.SheetWriter) (*response.DataResponse[string], error) {
	err := client.tcp_client.Send(sheet_writer)

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_data(sheet_writer)
	}

	response, err := client.get_str_response()

	if err != nil {
		if err := client.reconnect(); err != nil {
			return nil, err
		}

		return client.process_data(sheet_writer)
	}

	client.reconnect_attempts = 0
	return response, nil
}

func (client *PaperClient) process_has(sheet_writer *sheet_writer.SheetWriter) (*response.DataResponse[bool], error) {
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

func (client *PaperClient) process_size(sheet_writer *sheet_writer.SheetWriter) (*response.DataResponse[uint32], error) {
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

func (client *PaperClient) process_stats(sheet_writer *sheet_writer.SheetWriter) (*response.DataResponse[response.StatsData], error) {
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

func (client *PaperClient) get_response() (*response.Response, error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	if !is_ok {
		error, err := error_from_sheet(sheet_reader)

		if err != nil {
			return nil, err
		}

		return response.New(is_ok, &error), nil
	}

	return response.New(is_ok, nil), nil
}

func (client *PaperClient) get_str_response() (*response.DataResponse[string], error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	if !is_ok {
		error, err := error_from_sheet(sheet_reader)

		if err != nil {
			return nil, err
		}

		return response.NewData[string](is_ok, nil, &error), nil
	}

	data, err := sheet_reader.ReadString()

	if err != nil {
		return nil, err
	}

	return response.NewData(is_ok, &data, nil), nil
}

func (client *PaperClient) get_has_response() (*response.DataResponse[bool], error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	if !is_ok {
		error, err := error_from_sheet(sheet_reader)

		if err != nil {
			return nil, err
		}

		return response.NewData[bool](is_ok, nil, &error), nil
	}

	data, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	return response.NewData(is_ok, &data, nil), nil
}


func (client *PaperClient) get_size_response() (*response.DataResponse[uint32], error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	if !is_ok {
		error, err := error_from_sheet(sheet_reader)

		if err != nil {
			return nil, err
		}

		return response.NewData[uint32](is_ok, nil, &error), nil
	}

	data, err := sheet_reader.ReadU32()

	if err != nil {
		return nil, err
	}

	return response.NewData(is_ok, &data, nil), nil
}

func (client *PaperClient) get_stats_response() (*response.DataResponse[response.StatsData], error) {
	sheet_reader := sheet_reader.New(client.tcp_client)
	is_ok, err := sheet_reader.ReadBool()

	if err != nil {
		return nil, err
	}

	if is_ok {
		max_size, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		used_size, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		num_objects, err := sheet_reader.ReadU64()

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

		num_policies, err := sheet_reader.ReadU32()

		if err != nil {
			return nil, err
		}

		var policies []string

		for i := uint32(0); i < num_policies; i++ {
			policy, err := sheet_reader.ReadString()

			if err != nil {
				return nil, err
			}

			policies = append(policies, policy)
		}

		policy, err := sheet_reader.ReadString()

		if err != nil {
			return nil, err
		}

		is_auto_policy, err := sheet_reader.ReadBool()

		if err != nil {
			return nil, err
		}

		uptime, err := sheet_reader.ReadU64()

		if err != nil {
			return nil, err
		}

		stats := response.NewStatsData(
			max_size,
			used_size,
			num_objects,

			total_gets,
			total_sets,
			total_dels,

			miss_ratio,

			policies,
			policy,
			is_auto_policy,

			uptime,
		)

		return response.NewData(is_ok, &stats, nil), nil
	}

	error, err := error_from_sheet(sheet_reader)

	if err != nil {
		return nil, err
	}

	return response.NewData[response.StatsData](is_ok, nil, &error), nil
}

func error_from_sheet(sheet_reader *sheet_reader.SheetReader) (uint8, error) {
	code, err := sheet_reader.ReadU8()

	if err != nil {
		return PAPER_ERROR_INTERNAL, err
	}

	if code == 0 {
		cache_code, err := sheet_reader.ReadU8()

		if err != nil {
			return PAPER_ERROR_INTERNAL, err
		}

		switch (cache_code) {
			case 1: return PAPER_ERROR_KEY_NOT_FOUND, nil

			case 2: return PAPER_ERROR_ZERO_VALUE_SIZE, nil
			case 3: return PAPER_ERROR_EXCEEDING_VALUE_SIZE, nil

			case 4: return PAPER_ERROR_ZERO_CACHE_SIZE, nil

			case 5: return PAPER_ERROR_UNCONFIGURED_POLICY, nil
			case 6: return PAPER_ERROR_INVALID_POLICY, nil

			default: return PAPER_ERROR_INTERNAL, nil
		}
	}

	switch (code) {
		case 2: return PAPER_ERROR_MAX_CONNECTIONS_EXCEEDED, nil
		case 3: return PAPER_ERROR_UNAUTHORIZED, nil

		default: return PAPER_ERROR_INTERNAL, nil
	}
}

func parse_paper_addr(paper_addr string) (*string, error) {
	if !strings.HasPrefix(paper_addr, "paper://") {
		return nil, errors.New("Invalid paper address.");
	}

	addr := strings.Replace(paper_addr, "paper://", "", 1)

	return &addr, nil
}
