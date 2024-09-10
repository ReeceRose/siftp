package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/reecerose/siftp/utils"
)

func main() {
	listen, err := net.Listen(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		println("Failed to create TCP client", err)
		return
	}
	defer listen.Close()

	println("Server has started on PORT " + utils.SERVER_PORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			println("Failed to create accept new connection", err)
		}
		go recieveFile(conn)
	}
}

func recieveFile(conn net.Conn) {
	defer conn.Close()
	println("Received a request: " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	decoder := gob.NewDecoder(reader)
	var header utils.FileHeader
	err := decoder.Decode(&header)
	if err != nil {
		println("Failed to decode header", err)
		return
	}

	println("Protocol Version: %s", header.Version)
	println("Receiving file: %s", header.FileName)
	println("File Size: %d", header.FileSize)
	println("Checksum: %s", header.Checksum)

	os.Mkdir("./uploads", os.ModePerm)
	file, err := os.Create(fmt.Sprintf("./uploads/%s", header.FileName))
	if err != nil {
		println("Failed to create file to download to", err)
		return
	}
	defer file.Close()

	hash := sha256.New()
	limitedReader := io.LimitReader(reader, header.FileSize)
	writerToFile := io.MultiWriter(file, hash)

	bytesCopied, err := io.Copy(writerToFile, limitedReader)
	if err != nil {
		println("Failed to receive file data", err)
		return
	}

	calculatedChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	var serverResponse string
	if calculatedChecksum == header.Checksum {
		serverResponse = "File received successfully. Checksum matches.\n"
		println("File received successfully and checksum matches. Bytes transferred:", bytesCopied)
	} else {
		serverResponse = "Checksum mismatch!\n"
		println("Checksum mismatch!")
		println("Expected:", header.Checksum)
		println("Got:     ", calculatedChecksum)
	}

	writer.WriteString(serverResponse)
	writer.Flush()
}
