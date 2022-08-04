package main

import (
	"fmt"
	"log"
	"net"
)

const (
	CMD_SEND            = "send"
	CMD_DELIVER         = "deliver"
	CMD_SUSCRIBE        = "suscribe"
	DEFAULT_BUFFER_SIZE = 100 * 1024
)

type server struct {
	channels map[string]channel
	// connectedClients []client
}

// TODO: Compress the files
func (s *server) startServer(port int) {
	s.channels = make(map[string]channel)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	fmt.Printf("listening on port %d\n", port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.newClient(conn)
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Println("Client connected from", conn.RemoteAddr())
	cmdChan := make(chan command)
	contentChan := make(chan []byte)
	writeChan := make(chan []byte)
	newClient := client{conn: conn, writeChan: writeChan}
	go func(ch <-chan command) {
		for {
			cmd := <-ch
			s.handleCommand(&newClient, cmd, contentChan)
		}
	}(cmdChan)
	go newClient.startWriter()
	newClient.readRequest(cmdChan, contentChan)
}

func (s *server) handleCommand(client *client, cmd command, contentChan <-chan []byte) {

	switch cmd.Method {
	case CMD_SUSCRIBE:
		s.handleSuscribe(client, cmd)
	case CMD_SEND:
		s.handleSend(client, cmd, contentChan)
	}
}

func (s *server) handleSuscribe(suscriber *client, cmd command) {
	for _, cn := range cmd.Channels {
		chn, found := s.channels[cn]
		if found {
			chn.addClient(suscriber)
		} else {
			newChn := channel{name: cn, suscribedClients: map[string]*client{suscriber.conn.RemoteAddr().String(): suscriber}}
			s.channels[cn] = newChn
			log.Printf("New channel: %s", cn)
		}
		log.Printf("Client %s just suscribed to channel %s", (*suscriber).conn.RemoteAddr().String(), cn)
	}
}

func (s *server) handleSend(sender *client, cmd command, contentChan <-chan []byte) {
	senderAddress := sender.conn.RemoteAddr().String()
	cmd.Meta.SenderAddress = senderAddress
	// contentChans := make([]chan []byte, len(cmd.Channels))
	var contentChans []chan []byte

	for _, destChannel := range cmd.Channels {
		chn, found := s.channels[destChannel]
		if !found {
			//TODO: Add functionality to inform the client that channel doesn't exist'
			continue
		}
		deliverCmd := command{
			Method:   CMD_DELIVER,
			Meta:     metaData{SenderAddress: sender.conn.RemoteAddr().String()},
			Channels: []string{destChannel},
			FileInfo: cmd.FileInfo,
		}
		channelContentChan := make(chan []byte)
		go chn.broadcast(deliverCmd, channelContentChan)
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
