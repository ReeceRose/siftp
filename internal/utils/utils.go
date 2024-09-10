package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

const (
	TCP            = "tcp"
	SERVER_IP      = "0.0.0.0"
	SERVER_PORT    = "4567"
	SERVER_ADDRESS = SERVER_IP + ":" + SERVER_PORT // fmt.Sprintf doesn't work here
	VERSION_ONE    = "1.0"
)

func CalculateChecksum(file *os.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
