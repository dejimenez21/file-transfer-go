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
	Channels       map[string]*Channel
	channelsLock   sync.RWMutex
	requestCounter int64
}

func newServer() *Server {
	return &Server{
		Channels: make(map[string]*Channel),
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
			case req := <-client.RequestChan:
				s.handleRequest(client, req)
			case cl := <-client.Disconnect:
				s.disconnectClient(cl)
				return
			}

		}
	}(client)
	client.ReadRequest()
}

func (s *Server) handleRequest(client *Client, req models.Request) {

	switch req.Method {
	case models.REQ_SUSCRIBE:
		s.handleSuscriptionRequest(client, req)
	case models.REQ_SEND:
		s.handleFileSendingRequest(client, req)
	}
}

func (s *Server) handleSuscriptionRequest(suscriber *Client, cmd models.Request) {
	clientAddress := suscriber.Conn.RemoteAddr().String()
	s.channelsLock.Lock()
	for _, channelName := range cmd.Channels {
		channel, found := s.Channels[channelName]
		if found {
			channel.AddClient(suscriber)
		} else {
			NewChannel := newChannel(channelName)
			NewChannel.SuscribedClients[clientAddress] = suscriber
			s.Channels[channelName] = NewChannel
			log.Printf("New channel: %s", channelName)
		}
		log.Printf("Client %s suscribed to channel %s", clientAddress, channelName)
	}
	s.channelsLock.Unlock()
}

func (s *Server) handleFileSendingRequest(sender *Client, req models.Request) {
	senderAddress := sender.Conn.RemoteAddr().String()
	req.Meta.SenderAddress = senderAddress

	contentChans := s.beginBroadcasting(req, senderAddress)

	if contentChans == nil {
		s.sendAbortRequest(sender, "channels don't exist")
		for {
			select {
			case <-sender.ContentChan:
				continue
			case sender.AbortChan <- true:
				return
			}
		}
	}

	for i := 0; i < int(req.FileInfo.Size); {
		content := <-sender.ContentChan
		i += len(content)
		for _, channel := range contentChans {
			channel <- content
		}
	}

}

func (s *Server) beginBroadcasting(request models.Request, senderAddress string) []chan []byte {
	var contentChans []chan []byte

	for _, destChannel := range request.Channels {
		s.channelsLock.RLock()
		chn, found := s.Channels[destChannel]
		s.channelsLock.RUnlock()
		if !found {
			continue
		}
		deliverReq := models.NewDeliverRequest(senderAddress, s.newRequestId(), destChannel, request.FileInfo)
		channelContentChan := make(chan []byte)
		go chn.Broadcast(*deliverReq, channelContentChan)
		contentChans = append(contentChans, channelContentChan)
	}

	return contentChans
}

func (s *Server) newRequestId() int64 {
	s.requestCounter++
	return s.requestCounter
}

func (s *Server) disconnectClient(c *Client) {
	s.channelsLock.Lock()
	for key, channel := range s.Channels {
		channel.RemoveClient(c)
		if len(channel.SuscribedClients) < 1 {
			delete(s.Channels, key)
		}
	}
	s.channelsLock.Unlock()

}

func (s *Server) sendAbortRequest(client *Client, msg string) {
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
