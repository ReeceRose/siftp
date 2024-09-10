package types

import (
	"bufio"
	"encoding/gob"
	"io"
)

type FileHeader struct {
	Version  string
	FileName string
	FileSize int64
	Checksum string
}

func NewFileHeader(version string, fileName string, fileSize int64, checksum string) FileHeader {
	return FileHeader{
		Version:  version,
		FileName: fileName,
		FileSize: fileSize,
		Checksum: checksum,
	}
}

// Encode a header, requires a writer.Flush call
func (header FileHeader) Encode(writer io.Writer) error {
	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(header)
	if err != nil {
		return err
	}

	return nil
}

func DecodeHeader(reader *bufio.Reader) (*FileHeader, error) {
	decoder := gob.NewDecoder(reader)
	var header FileHeader
	err := decoder.Decode(&header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
