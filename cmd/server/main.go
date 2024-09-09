package main

import (
	"errors"
	"io"
	"net"
	"os"
	"time"

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
		go recieveFile(conn)
	}
}

func recieveFile(conn net.Conn) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	println("Received a request: " + conn.RemoteAddr().String())
	headerBuffer := make([]byte, 1024)

	_, err := conn.Read(headerBuffer)
	if err != nil {
		if errors.Is(err, io.EOF) {
			println("EOF recieved while trying to read header")
			conn.Write([]byte("Failed to read header"))
			return
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			println("DeadlineExceeded while trying to read header")
			conn.Write([]byte("Failed to read header"))
			return
		}
		panic(err)
	}

	conn.Write([]byte("File upload recieved"))
	println(string(headerBuffer))
}
