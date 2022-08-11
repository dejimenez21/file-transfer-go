package main

type cftpMessage interface {
	getMessageType() string
	process() error
}
