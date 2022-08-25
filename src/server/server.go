package main

import (
	"fmt"
	"log"
	"net"
	"server/cftp"
	"server/cftp/models"
	"sync"
)

const (
	DEFAULT_BUFFER_SIZE = 1024
)

type Server struct {
	channels       map[string]*channel
	channelsLock   sync.RWMutex
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
	client.ReadRequest()
}

func (s *Server) handleCommand(client *Client, cmd models.Request, contentChan <-chan []byte) {

	switch cmd.Method {
	case models.REQ_SUSCRIBE:
		s.handleSuscribe(client, cmd)
	case models.REQ_SEND:
		s.handleSend(client, cmd, contentChan)
	}
}

func (s *Server) handleSuscribe(suscriber *Client, cmd models.Request) {
	s.channelsLock.Lock()
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
	s.channelsLock.Unlock()
}

func (s *Server) handleSend(sender *Client, cmd models.Request, contentChan <-chan []byte) {
	senderAddress := sender.Conn.RemoteAddr().String()
	cmd.Meta.SenderAddress = senderAddress
	var contentChans []chan []byte

	for _, destChannel := range cmd.Channels {
		s.channelsLock.RLock()
		chn, found := s.channels[destChannel]
		s.channelsLock.RUnlock()
		if !found {
			continue
		}
		deliverCmd := models.Request{
			Method:   models.REQ_DELIVER,
			Meta:     models.MetaData{SenderAddress: sender.Conn.RemoteAddr().String(), RequestId: s.newRequestId()},
			Channels: []string{destChannel},
			FileInfo: cmd.FileInfo,
		}
		channelContentChan := make(chan []byte)
		go chn.Broadcast(deliverCmd, channelContentChan)
		contentChans = append(contentChans, channelContentChan)
	}

	if contentChans == nil {
		sendAbortRequest(sender, "channels don't exist")
		return
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
	s.channelsLock.Lock()
	for key, channel := range s.channels {
		channel.RemoveClient(c)
		if len(channel.suscribedClients) < 1 {
			delete(s.channels, key)
		}
	}
	s.channelsLock.Unlock()

}

func sendAbortRequest(client *Client, msg string) {
	req := models.NewAbortRequest(msg)
	reqBytes, err := cftp.SerializeRequest(*req)
	if err != nil {
		log.Printf("error serializing abort request: %v", err)
	}

	err = client.Write(reqBytes)
	if err != nil {
		log.Printf("error sending abort request: %v", err)
	}
}
