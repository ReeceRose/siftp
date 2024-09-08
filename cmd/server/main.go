package main

import (
	"net"

	"github.com/reecerose/siftp/utils"
)

func main() {
	listen, err := net.Listen(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	println("Server has started on PORT " + utils.SERVER_PORT)

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		println("New client recieved")

		conn.Close()
	}
}
