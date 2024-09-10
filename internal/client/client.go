package client

import (
	"net"

	"github.com/reecerose/siftp/internal/utils"
)

func SetupTCPClientConnection() (*net.TCPConn, error) {
	tcpServer, err := net.ResolveTCPAddr(utils.TCP, utils.SERVER_ADDRESS)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(utils.TCP, nil, tcpServer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
