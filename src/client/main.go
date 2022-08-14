package main

import (
	"flag"
	"log"
	"os"
	"time"
)

const (
	DEFAULT_RECEIVE_FOLDER_PATH      = "C:/Users/dejim/Documents/CFTP-Client/ReceiverFolder/"
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
	fact = new(factory)
	// client    = new(tcpClient)
	startTime  = time.Now().UnixNano()
	serverAddr string
)

func main() {
	fact.fileBroker = &fsBroker{}
	client := fact.getTcpClient()

	var channels arrayFlags

	receiveSet := flag.NewFlagSet("receive", flag.ExitOnError)
	//receiveDetached := receiveCmd.Bool("async", false, "Indicates if the receive operation should run asynchronously.")
	receivePath := receiveSet.String("path", DEFAULT_RECEIVE_FOLDER_PATH, "Folder where the received files will be stored.")
	receiveSet.Var(&channels, "ch", "Channel to receive files from.")
	receiveServerAddr := receiveSet.String("server", DEFAULT_SERVER_ADDR, "Address fo the CFTP server.")

	sendSet := flag.NewFlagSet("send", flag.ExitOnError)
	sendSet.Var(&channels, "ch", "Channel to send files to.")
	sendServerAddr := sendSet.String("server", DEFAULT_SERVER_ADDR, "Address fo the CFTP server.")

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

	if client.conn != nil {
		client.conn.Close()
	}
}

func handleReceiveCommand(cmd receiveCmd) {
	client := fact.getTcpClient()
	if len(cmd.channels) < 1 {
		log.Fatal("you need to provide at least one channel")
	}
	req := request{
		Method:   REQ_SUSCRIBE,
		Channels: cmd.channels,
	}

	client.sendRequest(req)
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
		log.Fatal("you need to specify the file to send")
	}
	if len(cmd.channels) < 1 {
		log.Fatal("you need to provide at least one channel")
	}

	// fileSize, err := fileBroker.getFileSize(cmd.filePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	contentChan := make(chan []byte)
	fInfo, err := fileBroker.loadFile(cmd.filePath, contentChan)
	if err != nil {
		log.Fatal(err)
	}
	req := request{
		Method:   REQ_SEND,
		Channels: cmd.channels,
		FileInfo: fInfo,
	}
	client.sendRequest(req)
	for content := range contentChan {
		client.sendFileContent(content)
	}
}
