package utils

const (
	TCP            = "tcp"
	SERVER_IP      = "0.0.0.0"
	SERVER_PORT    = "4567"
	SERVER_ADDRESS = SERVER_IP + ":" + SERVER_PORT // fmt.Sprintf doesn't work here
	TEST_FILE_PATH = "./test.txt"
	VERSION_ONE    = "1.0"
)

type FileHeader struct {
	Version  string
	FileName string
	FileSize int64
	Checksum string
}
