package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
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
	REQ_ABORT                        = "abort"
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
	log.SetFlags(0)
	fact.fileBroker = &fsBroker{}
	client := fact.getTcpClient()
	defer closeConnection(client)

	app := &cli.App{
		Name:  "client",
		Usage: "A CFTP client application",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "server", Usage: "Server address", Aliases: []string{"s"}, Value: DEFAULT_SERVER_ADDR},
		},
		Commands: []*cli.Command{
			{
				Name:  "send",
				Usage: "Send a file through several channels",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "channels", Aliases: []string{"ch"}, Usage: "Channels through which the file will be sent", Required: true},
				},
				Action: func(ctx *cli.Context) error {
					cmd := sendCmd{channels: ctx.StringSlice("channels")}
					cmd.filePath = ctx.Args().First()
					if cmd.filePath == "" {
						return fmt.Errorf("you need to specify the file to send")
					}
					serverAddr = ctx.String("server")
					handleSendCommand(cmd)
					return nil
				},
			},
			{
				Name:  "receive",
				Usage: "Start receiving files",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "channels", Aliases: []string{"ch"}, Usage: "Channels to suscribe", Required: true},
					&cli.PathFlag{Name: "outputDir", Aliases: []string{"o"}, Usage: "Save the received files to `DIRECTORY/`", Value: DEFAULT_RECEIVE_FOLDER_PATH},
				},
				Action: func(ctx *cli.Context) error {
					receivePath := ctx.Path("outputDir")
					fact.fileBroker.path = receivePath
					serverAddr = ctx.String("server")
					cmd := receiveCmd{channels: ctx.StringSlice("channels"), folderPath: receivePath}
					handleReceiveCommand(cmd)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func closeConnection(client *tcpClient) {
	if client.conn != nil {
		client.conn.Close()
	}
}

func handleReceiveCommand(cmd receiveCmd) {
	client := fact.getTcpClient()

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
	contentChan := make(chan []byte)
	abortChan := make(chan string)

	fInfo, err := fileBroker.loadFile(cmd.filePath, contentChan, abortChan)
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
	log.Printf("sending file %s ...", req.FileInfo.Name)
	go listenForAbortRequest(client, abortChan)

	for content := range contentChan {
		err = client.sendFileContent(content)
		if err != nil {
			log.Printf("error sending file content: %v\n", err)
			return
		}
	}
}

func listenForAbortRequest(client *tcpClient, abortChan chan<- string) {
	for {
		input, err := client.readInput()
		if err != nil {
			log.Println(err)
			continue
		}
		if input.getMessageType() != MSG_TYPE_REQ {
			log.Println("unexpected message type received: ", input.getMessageType())
			continue
		}

		req := input.(*request)
		if req.Method != REQ_ABORT {
			log.Println("unexpected request received: ", req.Method)
			continue
		}
		abortChan <- req.Meta.Message
		break
	}
}
