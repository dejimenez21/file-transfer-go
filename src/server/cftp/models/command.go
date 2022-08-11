package models

type Command struct {
	Method   string
	Meta     MetaData
	Channels []string
	FileInfo File
}

type MetaData struct {
	SenderAddress string
	RequestId     int
}
