package types

type FileHeader struct {
	Version  string
	FileName string
	FileSize int64
	Checksum string
}
