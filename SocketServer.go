package go_server_core

import (
	"log"
	"net"
)

// TCP,UDP,HTTP Server
type SocketServer struct {
	serverType               string // TCP, UDP
	port                     string
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, pktName string, message string)
}

// Creates new tcp SocketServer instance
func NewSocketServer(serverType string, port string) *SocketServer {
	log.Println("Creating ", serverType, " SocketServer with port", port)
	SocketServer := &SocketServer{
		serverType: serverType,
		port:       port,
	}

	SocketServer.OnNewClient(func(c *Client) {})
	SocketServer.OnNewMessage(func(c *Client, pktName string, message string) {})
	SocketServer.OnClientConnectionClosed(func(c *Client, err error) {})

	go SocketServer.Listen()

	return SocketServer
}

// Listen starts network SocketServer
func (s *SocketServer) Listen() {
	if s.serverType == "tcp" {
		listener, err := net.Listen(s.serverType, s.port)
		if err != nil {
			log.Fatal("Error starting ", s.serverType, " SocketServer.\r\n", err)
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

// Called right after SocketServer starts listening new client
func (s *SocketServer) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// Called right after connection closed
func (s *SocketServer) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// Called when Client receives new message
func (s *SocketServer) OnNewMessage(callback func(c *Client, pktName string, message string)) {
	s.onNewMessage = callback
}
