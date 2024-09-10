package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/reecerose/siftp/internal/logging"
	"github.com/reecerose/siftp/internal/types"
	"github.com/reecerose/siftp/internal/utils"
)

func main() {
	logFile, err := logging.SetupLogging("client.log")
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	listen, err := net.Listen(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		log.Println("Failed to create TCP client", err)
		return
	}
	defer listen.Close()

	log.Println("Server has started on PORT " + utils.SERVER_PORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Failed to create accept new connection", err)
		}
		go recieveFile(conn)
	}
}

func recieveFile(conn net.Conn) {
	defer conn.Close()
	log.Println("Received a request: " + conn.RemoteAddr().String())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	decoder := gob.NewDecoder(reader)
	var header types.FileHeader
	err := decoder.Decode(&header)
	if err != nil {
		log.Println("Failed to decode header", err)
		return
	}

	log.Printf("Protocol Version: %s\n", header.Version)
	log.Printf("Receiving file: %s\n", header.FileName)
	log.Printf("File Size: %d\n", header.FileSize)
	log.Printf("Checksum: %s\n", header.Checksum)

	os.Mkdir("./uploads", os.ModePerm)
	file, err := os.Create(fmt.Sprintf("./uploads/%s", header.FileName))
	if err != nil {
		log.Println("Failed to create file to download to", err)
		return
	}
	defer file.Close()

	hash := sha256.New()
	limitedReader := io.LimitReader(reader, header.FileSize)
	writerToFile := io.MultiWriter(file, hash)

	bytesCopied, err := io.Copy(writerToFile, limitedReader)
	if err != nil {
		log.Println("Failed to receive file data", err)
		return
	}

	calculatedChecksum := fmt.Sprintf("%x", hash.Sum(nil))
	var serverResponse string
	if calculatedChecksum == header.Checksum {
		serverResponse = "File received successfully. Checksum matches.\n"
		log.Println("File received successfully and checksum matches. Bytes transferred:", bytesCopied)
	} else {
		serverResponse = "Checksum mismatch!\n"
		log.Println("Checksum mismatch!")
		log.Println("Expected:", header.Checksum)
		log.Println("Got:     ", calculatedChecksum)
	}

	writer.WriteString(serverResponse)
	writer.Flush()
}
