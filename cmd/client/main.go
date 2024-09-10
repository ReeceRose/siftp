package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"

	"github.com/reecerose/siftp/utils"
)

func main() {
	println("Client started")

	tcpServer, err := net.ResolveTCPAddr(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		println("Failed to create TCP client", err)
		return
	}

	conn, err := net.DialTCP(utils.TCP, nil, tcpServer)
	if err != nil {
		println("Failed to create TCP connection", err)
		return
	}
	defer conn.Close()

	println("Connected to " + utils.SERVER_ADDRESS)

	file, err := os.OpenFile(utils.TEST_FILE_PATH, 0, fs.FileMode(os.O_RDONLY))
	if err != nil {
		println("Failed to open file", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		println("Failed to get file info", err)
		return
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		println("Failed to calculate checksum", err)
		return
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	if _, err := file.Seek(0, 0); err != nil {
		println("Failed to reset file pointer", err)
		return
	}

	header := utils.FileHeader{
		Version:  utils.VERSION_ONE,
		FileName: fileInfo.Name(),
		FileSize: fileInfo.Size(),
		Checksum: checksum,
	}

	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(header)
	if err != nil {
		println("Failed to encode header", err)
		return
	}

	// Send File Data
	bytesCopied, err := io.Copy(writer, file)
	if err != nil {
		println("Failed to copy data to writer", err)
		return
	}

	// Flush the writer again after file transfer
	writer.Flush()

	println("File sent successfully, bytes transferred:", bytesCopied)

	response, err := reader.ReadString('\n')
	if err != nil {
		println("Failed to get response from server", err)
		return
	}

	println("Server response:", response)
}
