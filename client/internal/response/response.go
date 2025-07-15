/*
 * Copyright (c) Kia Shakiba
 *
 * This source code is licensed under the GNU AGPLv3 license found in the
 * LICENSE file in the root directory of this source tree.
 */

package response

type Response struct {
	is_ok bool
	error *uint8
}

type DataResponse[T any] struct {
	is_ok bool
	data *T
	error *uint8
}

type StatusData struct {
	pid uint32

	max_size uint64
	used_size uint64
	num_objects uint64

	rss uint64
	hwm uint64

	total_gets uint64
	total_sets uint64
	total_dels uint64

	miss_ratio float64

	policies []string
	policy string
	is_auto_policy bool

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

func NewStatusData(
	pid uint32,

	max_size uint64,
	used_size uint64,
	num_objects uint64,

	rss uint64,
	hwm uint64,

	total_gets uint64,
	total_sets uint64,
	total_dels uint64,

	miss_ratio float64,

	policies []string,
	policy string,
	is_auto_policy bool,

	uptime uint64,
) StatusData {
	return StatusData {
		pid,

		max_size,
		used_size,
		num_objects,

		rss,
		hwm,

		total_gets,
		total_sets,
		total_dels,

		miss_ratio,

		policies,
		policy,
		is_auto_policy,

		uptime,
	}
}

func (status StatusData) MaxSize() uint64 {
	return status.max_size
}

func (status StatusData) UsedSize() uint64 {
	return status.used_size
}

func (status StatusData) NumObjects() uint64 {
	return status.num_objects
}

func (status StatusData) TotalGets() uint64 {
	return status.total_gets
}

func (status StatusData) TotalSets() uint64 {
	return status.total_sets
}

func (status StatusData) TotalDels() uint64 {
	return status.total_dels
}

func (status StatusData) MissRatio() float64 {
	return status.miss_ratio
}

func (status StatusData) Policies() []string {
	return status.policies
}

func (status StatusData) Policy() string {
	return status.policy
}

func (status StatusData) IsAutoPolicy() bool {
	return status.is_auto_policy
}

func (status StatusData) Uptime() uint64 {
	return status.uptime
}
