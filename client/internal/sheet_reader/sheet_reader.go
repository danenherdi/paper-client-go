package sheet_reader

import (
	"math"
	"encoding/binary"
	"internal/tcp_client"
)

type SheetReader struct {
	tcp_client *tcp_client.TcpClient
}

func New(tcp_client *tcp_client.TcpClient) *SheetReader {
	return &SheetReader {
		tcp_client,
	}
}

func (sheet *SheetReader) ReadU8() (uint8, error) {
	data := make([]byte, 1)
	_, err := sheet.tcp_client.GetConn().Read(data)

	if err != nil {
		return 0, err
	}

	return data[0], nil
}

func (sheet *SheetReader) ReadU32() (uint32, error) {
	data := make([]byte, 4)
	_, err := sheet.tcp_client.GetConn().Read(data)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(data), nil
}

func (sheet *SheetReader) ReadU64() (uint64, error) {
	data := make([]byte, 8)
	_, err := sheet.tcp_client.GetConn().Read(data)

	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(data), nil
}

func (sheet *SheetReader) ReadF64() (float64, error) {
	data := make([]byte, 8)
	_, err := sheet.tcp_client.GetConn().Read(data)

	if err != nil {
		return 0, err
	}

	bits := binary.LittleEndian.Uint64(data)
	return math.Float64frombits(bits), nil
}

func (sheet *SheetReader) ReadBool() (bool, error) {
	data, err := sheet.ReadU8()

	if err != nil {
		return false, err
	}

	return data == '!', nil
}

func (sheet *SheetReader) ReadString() (string, error) {
	length, err := sheet.ReadU32()

	if err != nil {
		return "", err
	}

	data := make([]byte, length)
	_, err = sheet.tcp_client.GetConn().Read(data)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
