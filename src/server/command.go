package main

type command struct {
	Method   string
	Meta     metaData
	Channels []string
	FileInfo fileMeta
	sender   *client
}

type metaData struct {
	HasFile bool
}
