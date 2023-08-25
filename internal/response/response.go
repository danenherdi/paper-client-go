package response

type Response[T any] struct {
	is_ok bool
	data T
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

func New[T any](is_ok bool, data T) *Response[T] {
	return &Response[T] {
		is_ok,
		data,
	}
}

func (response *Response[T]) IsOk() bool {
	return response.is_ok
}

func (response *Response[T]) Data() T {
	return response.data
}

func NewStatsResponse(
	max_size uint64,
	used_size uint64,

	total_gets uint64,
	total_sets uint64,
	total_dels uint64,

	miss_ratio float64,

	policy uint8,
	uptime uint64,
) *Response[StatsData] {
	return New(
		true,

		StatsData {
			max_size,
			used_size,

			total_gets,
			total_sets,
			total_dels,

			miss_ratio,

			policy,
			uptime,
		},
	)
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
