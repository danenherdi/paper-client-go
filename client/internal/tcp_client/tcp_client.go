package tcp_client

import (
	"net"
	"errors"
	"strconv"
	"internal/sheet_writer"
)

type TcpClient struct {
	conn *net.TCPConn
}

func Connect(host string, port uint32) (*TcpClient, error) {
	server, err := net.ResolveTCPAddr("tcp", host + ":" + strconv.Itoa(int(port)))

	if err != nil {
		return nil, errors.New("Invalid host or port.")
	}

	conn, err := net.DialTCP("tcp", nil, server)

	if err != nil {
		return nil, errors.New("Could not connect to server.")
	}

	client := TcpClient {
		conn,
	}

	return &client, nil
}

func (client *TcpClient) GetConn() *net.TCPConn {
	return client.conn
}

func (client *TcpClient) Send(sheet *sheet_writer.SheetWriter) (error) {
	_, err := client.conn.Write(sheet.GetBuf())
	return err
}
