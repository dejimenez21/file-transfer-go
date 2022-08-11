package main

import "fmt"

type delivery struct {
	DeliveryId int
	Seq        int
	Size       int
	Content    []byte
}

func (del *delivery) getMessageType() string {
	return "chunk"
}

func (chunk *delivery) process() (err error) {
	fileBroker := fact.getFileBroker()
	err = fileBroker.saveChunk(chunk)
	if err != nil {
		err = fmt.Errorf("chunk disposed: %v", err)
	}
	return err
}
