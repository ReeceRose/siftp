package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"

	"github.com/reecerose/siftp/utils"
)

func main() {
	logFile, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if len(os.Args) < 2 {
		log.Panicln("Usage: go run cmd/client/main.go <file_path>")
		return
	}

	filePath := os.Args[1]
	log.Println("Client started")

	tcpServer, err := net.ResolveTCPAddr(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		log.Println("Failed to create TCP client", err)
		return
	}

	conn, err := net.DialTCP(utils.TCP, nil, tcpServer)
	if err != nil {
		log.Println("Failed to create TCP connection", err)
		return
	}
	defer conn.Close()

	log.Println("Connected to " + utils.SERVER_ADDRESS)

	file, err := os.OpenFile(filePath, 0, fs.FileMode(os.O_RDONLY))
	if err != nil {
		log.Println("Failed to open file", err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Failed to get file info", err)
		return
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Println("Failed to calculate checksum", err)
		return
	}
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Failed to reset file pointer", err)
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
		log.Println("Failed to encode header", err)
		return
	}

	// Send File Data
	bytesCopied, err := io.Copy(writer, file)
	if err != nil {
		log.Println("Failed to copy data to writer", err)
		return
	}

	// Flush the writer again after file transfer
	writer.Flush()

	log.Println("File sent successfully, bytes transferred:", bytesCopied)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to get response from server", err)
		return
	}

	log.Println("Server response:", response)
}
