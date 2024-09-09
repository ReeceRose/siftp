package main

import (
	"errors"
	"fmt"
	"io"
	"net"

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

	// TODO: file uploadw

	uploadMessage := make([]byte, 1024)

	for {
		n, err := conn.Read(uploadMessage)

		if err != nil {
			if errors.Is(err, io.EOF) {
				println("Connection closed, upload failed")
				println(string(uploadMessage[:n]))
				break
			}

			panic(err)
		}

		// If bytes were read, process the data
		if n > 0 {
			println(string(uploadMessage[:n]))
			break
		}
	}
}
