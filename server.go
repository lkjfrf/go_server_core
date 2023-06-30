package go_server_core

import (
	"log"
	"net"
)

// TCP,UDP,HTTP Server
type server struct {
	serverType               string // TCP, UDP
	port                     string // Address to open connection: localhost:9999
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, pktName string, message string)
}

// Creates new tcp server instance
func New(serverType string, port string) *server {
	log.Println("Creating ", serverType, " server with port", port)
	server := &server{
		serverType: serverType,
		port:       port,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, pktName string, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// Listen starts network server
func (s *server) Listen() {
	if s.serverType == "tcp" {
		listener, err := net.Listen(s.serverType, s.port)
		if err != nil {
			log.Fatal("Error starting ", s.serverType, " server.\r\n", err)
		}
		defer listener.Close()

		for {
			conn, _ := listener.Accept()
			client := &Client{
				ServerType: s.serverType,
				Tcpconn:    conn,
				Server:     s,
			}
			go client.listen()
		}
	} else {
		addr := net.UDPAddr{
			Port: 9999,
			IP:   net.ParseIP("localhost"),
		}
		udpServer, err := net.ListenUDP("udp", &addr)
		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			ServerType: s.serverType,
			Udpconn:    udpServer,
			Server:     s,
		}
		go client.listen()
	}
}

//////////////////////////////////////
//			CALLBACKS				//
/////////////////////////////////////

// Called right after server starts listening new client
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *server) OnNewMessage(callback func(c *Client, pktName string, message string)) {
	s.onNewMessage = callback
}
