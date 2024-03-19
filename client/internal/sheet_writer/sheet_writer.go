package sheet_writer

import (
	"encoding/binary"
)

type SheetWriter struct {
	buf []byte
}

func New() *SheetWriter {
	return &SheetWriter {
		buf: []byte{},
	}
}

func (sheet *SheetWriter) GetBuf() []byte {
	return sheet.buf
}

func (sheet *SheetWriter) WriteU8(value uint8) {
	sheet.buf = append(sheet.buf, value)
}

func (sheet *SheetWriter) WriteU32(value uint32) {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, value)
	sheet.buf = append(sheet.buf, data...)
}

func (sheet *SheetWriter) WriteU64(value uint64) {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, value)
	sheet.buf = append(sheet.buf, data...)
}

func (sheet *SheetWriter) WriteString(value string) {
	length := len(value)
	sheet.WriteU32(uint32(length))
	sheet.buf = append(sheet.buf, value...)
}
