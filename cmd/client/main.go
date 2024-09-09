package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"

	"github.com/reecerose/siftp/utils"
)

func main() {
	fmt.Println("Client started")

	tcpServer, err := net.ResolveTCPAddr(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialTCP(utils.TCP, nil, tcpServer)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	println("Connected to " + utils.SERVER_ADDRESS)

	file, err := os.OpenFile(utils.TEST_FILE_PATH, 0, fs.FileMode(os.O_RDONLY))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// TODO: actual file upload
	conn.Write([]byte(file.Name()))
	if err != nil {
		panic(err)
	}

	uploadMessage := make([]byte, 1024)

	for {
		n, err := conn.Read(uploadMessage)

		if err != nil {
			println(string(uploadMessage[:n]))
			if errors.Is(err, io.EOF) {
				println("Connection closed")
				break
			}
			panic(err)
		}

		if n > 0 {
			println(string(uploadMessage[:n]))
		}
	}
}
