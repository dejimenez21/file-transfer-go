package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DEFAULT_RECEIVE_FOLDER_PATH = "C:/Users/dejim/Documents/CFTP-Client/ReceiverFolder/"
	MAX_FILE_SIZE               = 100 * 1024
	CMD_SEND                    = "send"
	CMD_RECEIVE                 = "receive"
	REQ_SUSCRIBE                = "suscribe"
	REQ_DELIVER                 = "deliver"
	REQ_SEND                    = "send"
)

var (
	client    = new(tcpClient)
	startTime = time.Now().UnixNano()
)

func main() {
	var channels arrayFlags

	receiveSet := flag.NewFlagSet("receive", flag.ExitOnError)
	//receiveDetached := receiveCmd.Bool("async", false, "Indicates if the receive operation should run asynchronously.")
	receivePath := receiveSet.String("path", DEFAULT_RECEIVE_FOLDER_PATH, "Folder where the received files will be stored.")
	receiveSet.Var(&channels, "ch", "Channel to receive files from.")

	sendSet := flag.NewFlagSet("send", flag.ExitOnError)
	sendSet.Var(&channels, "ch", "Channel to send files to.")

	method := os.Args[1]

	switch method {
	case CMD_RECEIVE:
		receiveSet.Parse(os.Args[2:len(os.Args)])
		cmd := receiveCmd{channels: channels, folderPath: *receivePath}
		handleReceiveCommand(cmd)
	case CMD_SEND:
		sendSet.Parse(os.Args[2 : len(os.Args)-1])
		cmd := sendCmd{channels: channels}
		cmd.filePath = os.Args[len(os.Args)-1]
		handleSendCommand(cmd)
	}

	if client.conn != nil {
		client.conn.Close()
	}
}

func handleReceiveCommand(cmd receiveCmd) {
	if len(cmd.channels) < 1 {
		log.Fatal("you need to provide at least one channel")
	}
	req := request{
		Method:   REQ_SUSCRIBE,
		Channels: cmd.channels,
	}
	client.sendRequest(req)
	fileChn := make(chan *request)
	go client.startReceiving(fileChn)

	for {
		req := <-fileChn
		switch req.Method {
		case REQ_DELIVER:
			fileBroker := fsBroker{path: cmd.folderPath}
			newFilePath, err := fileBroker.saveFile(req.FileInfo)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Printf("%s file received from %s throug channel %s saved as %s\n", req.FileInfo.FullName(), req.Meta.SenderAddress, req.Channels[0], newFilePath)
		}
	}
}

func handleSendCommand(cmd sendCmd) {
	if cmd.filePath == "" {
		log.Fatal("you need to specify the file to send")
	}
	if len(cmd.channels) < 1 {
		log.Fatal("you need to provide at least one channel")
	}
	fileBroker := fsBroker{}

	fileSize, err := fileBroker.getFileSize(cmd.filePath)
	if err != nil {
		log.Fatal(err)
	}

	if fileSize > MAX_FILE_SIZE {
		//TODO: add logic to handle bigger files
		log.Fatalf("file size is too large")
	}

	fInfo, err := fileBroker.loadFile(cmd.filePath)
	if err != nil {
		log.Fatal(err)
	}
	fInfo.Size = fileSize
	req := request{
		Method:   REQ_SEND,
		Meta:     metaData{HasFileContent: true},
		Channels: cmd.channels,
		FileInfo: fInfo,
	}

	client.sendRequest(req)
}

// func getFilePath() (path string, err error) {
// 	args := os.Args
// 	path = args[1]
// 	fmt.Println(path)
// 	return
// }

// func connect() (conn net.Conn, err error) {
// 	var d net.Dialer
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()

// 	conn, err = d.DialContext(ctx, "tcp", "localhost:8888")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return
// }
