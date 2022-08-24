package main

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	DEFAULT_RECEIVE_FOLDER_PATH      = "./ReceivedFiles/"
	MAX_FILE_SIZE                    = 1024
	EOT                         byte = 0x04
	CMD_SEND                         = "send"
	CMD_RECEIVE                      = "receive"
	REQ_SUSCRIBE                     = "suscribe"
	REQ_DELIVER                      = "deliver"
	REQ_SEND                         = "send"
	MSG_TYPE_REQ                     = "request"
	MSG_TYPE_CHUNK                   = "chunk"
	DEFAULT_SERVER_ADDR              = "localhost:8888"
)

var (
	fact       = new(factory)
	startTime  = time.Now().UnixNano()
	serverAddr string
)

func main() {
	fact.fileBroker = &fsBroker{}
	client := fact.getTcpClient()
	defer closeConnection(client)

	var channels arrayFlags

	receiveSet := flag.NewFlagSet("receive", flag.ExitOnError)
	receivePath := receiveSet.String("path", DEFAULT_RECEIVE_FOLDER_PATH, "Folder where the received files will be stored.")
	receiveSet.Var(&channels, "ch", "Channel to receive files from.")
	receiveServerAddr := receiveSet.String("server", DEFAULT_SERVER_ADDR, "Address for the CFTP server.")

	sendSet := flag.NewFlagSet("send", flag.ExitOnError)
	sendSet.Var(&channels, "ch", "Channel to send files to.")
	sendServerAddr := sendSet.String("server", DEFAULT_SERVER_ADDR, "Address for the CFTP server.")

	method := os.Args[1]

	switch method {
	case CMD_RECEIVE:
		receiveSet.Parse(os.Args[2:len(os.Args)])
		fact.fileBroker.path = *receivePath
		serverAddr = *receiveServerAddr
		cmd := receiveCmd{channels: channels, folderPath: *receivePath}
		handleReceiveCommand(cmd)
	case CMD_SEND:
		sendSet.Parse(os.Args[2 : len(os.Args)-1])
		cmd := sendCmd{channels: channels}
		cmd.filePath = os.Args[len(os.Args)-1]
		serverAddr = *sendServerAddr
		handleSendCommand(cmd)
	}

}

func closeConnection(client *tcpClient) {
	if client.conn != nil {
		client.conn.Close()
	}
}

func handleReceiveCommand(cmd receiveCmd) {
	client := fact.getTcpClient()
	if len(cmd.channels) < 1 {
		log.Print("you need to provide at least one channel")
		return
	}
	req := request{
		Method:   REQ_SUSCRIBE,
		Channels: cmd.channels,
	}

	err := client.sendRequest(req)
	if err != nil {
		log.Printf("error sending request: %v\n", err)
		return
	}
	for {
		input, err := client.readInput()
		if err != nil {
			log.Printf("error reading input: %v\n", err)
			continue
		}
		err = input.process()
		if err != nil {
			log.Println(err)
		}
	}
}

func handleSendCommand(cmd sendCmd) {
	client := fact.getTcpClient()
	fileBroker := fact.getFileBroker()

	if cmd.filePath == "" {
		log.Print("you need to specify the file to send")
		return
	}
	if len(cmd.channels) < 1 {
		log.Print("you need to provide at least one channel")
		return
	}

	contentChan := make(chan []byte)
	fInfo, err := fileBroker.loadFile(cmd.filePath, contentChan)
	if err != nil {
		log.Print(err)
		return
	}
	req := request{
		Method:   REQ_SEND,
		Channels: cmd.channels,
		FileInfo: fInfo,
	}
	err = client.sendRequest(req)
	if err != nil {
		log.Printf("error sending send request: %v\n", err)
		return
	}
	for content := range contentChan {
		err = client.sendFileContent(content)
		if err != nil {
			log.Printf("error sending file content: %v\n", err)
			return
		}
	}
}
