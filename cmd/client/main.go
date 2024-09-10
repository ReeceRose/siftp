package main

import (
	"bufio"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/reecerose/siftp/internal/client"
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

	if len(os.Args) < 2 {
		log.Panicln("Usage: go run cmd/client/main.go <file_path>")
		return
	}

	filePath := os.Args[1]
	log.Println("Client started")

	conn, err := client.SetupTCPClientConnection()
	if err != nil {
		log.Println("Failed to create TCP client", err)
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

	checksum, err := utils.CalculateChecksum(file)
	if err != nil {
		log.Println("Failed to calculate checksum", err)
		return
	}

	if _, err := file.Seek(0, 0); err != nil {
		log.Println("Failed to reset file pointer", err)
		return
	}

	header := types.NewFileHeader(utils.VERSION_ONE, file.Name(), fileInfo.Size(), checksum)
	err = header.Encode(writer)
	if err != nil {
		log.Println("Failed to encode header", err)
		return
	}
	writer.Flush()

	bytesCopied, err := io.Copy(writer, file)
	if err != nil {
		log.Println("Failed to copy data to writer", err)
		return
	}

	writer.Flush()

	log.Println("File sent successfully, bytes transferred:", bytesCopied)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Failed to get response from server", err)
		return
	}

	log.Println("Server response:", response)
}
