package response

const (
	ERROR_INTERNAL uint8 					= 0

	ERROR_UNREACHABLE_SERVER uint8			= 1
	ERROR_MAX_CONNECTIONS_EXCEEDED uint8	= 2
	ERROR_UNAUTHORIZED uint8				= 3

	ERROR_KEY_NOT_FOUND uint8				= 4

	ERROR_ZERO_VALUE_SIZE uint8				= 5
	ERROR_EXCEEDING_VALUE_SIZE uint8		= 6

	ERROR_ZERO_CACHE_SIZE uint8				= 7
)

type Response struct {
	is_ok bool
	error *uint8
}

type DataResponse[T any] struct {
	is_ok bool
	data *T
	error *uint8
}

type StatsData struct {
	max_size uint64
	used_size uint64

	total_gets uint64
	total_sets uint64
	total_dels uint64

	miss_ratio float64

	policy uint8
	uptime uint64
}

func New(is_ok bool, error *uint8) *Response {
	return &Response {
		is_ok,
		error,
	}
}

func NewData[T any](is_ok bool, data *T, error *uint8) *DataResponse[T] {
	return &DataResponse[T] {
		is_ok,
		data,
		error,
	}
}

func (response *Response) IsOk() bool {
	return response.is_ok
}

func (response *DataResponse[T]) IsOk() bool {
	return response.is_ok
}

func (response *DataResponse[T]) Data() *T {
	return response.data
}

func (response *Response) Error() uint8 {
	return *response.error
}


func (response *DataResponse[T]) Error() uint8 {
	return *response.error
}

func NewStatsData(
	max_size uint64,
	used_size uint64,

	total_gets uint64,
	total_sets uint64,
	total_dels uint64,

	miss_ratio float64,

	policy uint8,
	uptime uint64,
) StatsData {
	return StatsData {
		max_size,
		used_size,

		total_gets,
		total_sets,
		total_dels,

		miss_ratio,

		policy,
		uptime,
	}
}

func (stats StatsData) MaxSize() uint64 {
	return stats.max_size
}

func (stats StatsData) UsedSize() uint64 {
	return stats.used_size
}

func (stats StatsData) TotalGets() uint64 {
	return stats.total_gets
}

func (stats StatsData) TotalSets() uint64 {
	return stats.total_sets
}

func (stats StatsData) TotalDels() uint64 {
	return stats.total_dels
}

func (stats StatsData) MissRatio() float64 {
	return stats.miss_ratio
}

func (stats StatsData) Policy() uint8 {
	return stats.policy
}

func (stats StatsData) Uptime() uint64 {
	return stats.uptime
}
