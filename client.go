package go_server_core

import (
	"errors"
	"io"
	"log"
	"net"
	"syscall"
)

// Client holds info about connection
type Client struct {
	ServerType string //TCP or UDP
	Tcpconn    net.Conn
	Udpconn    net.PacketConn
	UdpAdr     net.Addr
	Server     *SocketServer
}

// Read client data from channel
func (c *Client) listen() {
	if c.ServerType == "tcp" {
		c.Server.onNewClientCallback(c)

		header := make([]byte, 4)
		for {
			_, err := c.Tcpconn.Read(header)
			namesize, datasize := ParseHeader(header)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, syscall.Errno(10054)) {
					c.Tcpconn.Close()
					c.Server.onClientConnectionClosed(c, err)
				}
				return
			}

			recvBuf := make([]byte, namesize+datasize)
			n, err := c.Tcpconn.Read(recvBuf)
			if err != nil {
				log.Println("TCP error", err)
				c.Tcpconn.Close()
				c.Server.onClientConnectionClosed(c, err)
				return
			}
			if 0 < n {

				pktname, jsondata := ExtractData(namesize, datasize, recvBuf)

				c.Server.onNewMessage(c, pktname, jsondata)
			}
		}

	} else {
		var err error
		recvBuf := make([]byte, 4096)
		var n int
		for {
			n, c.UdpAdr, err = c.Udpconn.ReadFrom(recvBuf)
			if err != nil {
				log.Println("UDP error", err)
				return
			}
			if 0 < n {
				//	c.Server.onNewMessage(c, recvBuf[:n])
			}
		}
	}
}

// Send bytes to client
func (c *Client) SendBytes(b []byte) error {
	if c.ServerType == "tcp" {
		_, err := c.Tcpconn.Write(b)
		if err != nil {
			c.Tcpconn.Close()
			c.Server.onClientConnectionClosed(c, err)
		}
		return err
	} else {
		_, err := c.Udpconn.WriteTo(b, c.UdpAdr)
		if err != nil {
			c.Udpconn.Close()
			c.Server.onClientConnectionClosed(c, err)
		}
		return err
	}
}

func (c *Client) Close() error {
	if c.ServerType == "tcp" {
		return c.Udpconn.Close()
	} else {
		return c.Tcpconn.Close()
	}
}
