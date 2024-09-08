package main

import (
	"fmt"
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

	println("Connected to " + utils.SERVER_ADDRESS)

	defer conn.Close()
}
