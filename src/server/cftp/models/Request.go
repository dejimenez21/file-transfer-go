package models

type Request struct {
	Method   string
	Meta     MetaData
	Channels []string
	FileInfo File
}

type MetaData struct {
	SenderAddress string
	RequestId     int64
}
