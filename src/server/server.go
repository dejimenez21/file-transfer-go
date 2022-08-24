package main

import (
	"fmt"
	"log"
	"net"
	"server/cftp/models"
)

const (
	CMD_SEND            = "send"
	CMD_DELIVER         = "deliver"
	CMD_SUSCRIBE        = "suscribe"
	DEFAULT_BUFFER_SIZE = 1024
)

type Server struct {
	channels       map[string]*channel
	requestCounter int64
}

func newServer() *Server {
	return &Server{
		channels: make(map[string]*channel),
	}
}

func (s *Server) StartServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Println("Client connected from", conn.RemoteAddr())

	client := newClient(conn)
	go func(client *Client) {
		for {
			select {
			case cmd := <-client.CmdChan:
				s.handleCommand(client, cmd, client.ContentChan)
			case cl := <-client.Disconnect:
				s.disconnectClient(cl)
			}

		}
	}(client)
	go client.startWriter()
	client.readRequest()
}

func (s *Server) handleCommand(client *Client, cmd models.Command, contentChan <-chan []byte) {

	switch cmd.Method {
	case CMD_SUSCRIBE:
		s.handleSuscribe(client, cmd)
	case CMD_SEND:
		s.handleSend(client, cmd, contentChan)
	}
}

func (s *Server) handleSuscribe(suscriber *Client, cmd models.Command) {
	for _, cn := range cmd.Channels {
		chn, found := s.channels[cn]
		if found {
			chn.AddClient(suscriber)
		} else {
			newChn := &channel{name: cn, suscribedClients: map[string]*Client{suscriber.Conn.RemoteAddr().String(): suscriber}}
			s.channels[cn] = newChn
			log.Printf("New channel: %s", cn)
		}
		log.Printf("Client %s just suscribed to channel %s", (*suscriber).Conn.RemoteAddr().String(), cn)
	}
}

func (s *Server) handleSend(sender *Client, cmd models.Command, contentChan <-chan []byte) {
	senderAddress := sender.Conn.RemoteAddr().String()
	cmd.Meta.SenderAddress = senderAddress
	var contentChans []chan []byte

	for _, destChannel := range cmd.Channels {
		chn, found := s.channels[destChannel]
		if !found {
			//TODO: Add functionality to inform the client that channel doesn't exist'
			continue
		}
		deliverCmd := models.Command{
			Method:   CMD_DELIVER,
			Meta:     models.MetaData{SenderAddress: sender.Conn.RemoteAddr().String(), RequestId: int(s.newRequestId())},
			Channels: []string{destChannel},
			FileInfo: cmd.FileInfo,
		}
		channelContentChan := make(chan []byte)
		go chn.Broadcast(deliverCmd, channelContentChan)
		contentChans = append(contentChans, channelContentChan)
	}
	for i := 0; i < int(cmd.FileInfo.Size); {
		content := <-contentChan
		i += len(content)
		for _, channel := range contentChans {
			channel <- content
		}
	}

}

func (s *Server) newRequestId() int64 {
	s.requestCounter++
	return s.requestCounter
}

func (s *Server) disconnectClient(c *Client) {
	for key, channel := range s.channels {
		channel.RemoveClient(c)
		if len(channel.suscribedClients) < 1 {
			delete(s.channels, key)
		}
	}
}
