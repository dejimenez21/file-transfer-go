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
	Message       string
}

func NewAbortRequest(msg string) *Request {
	return &Request{
		Method: REQ_ABORT,
		Meta:   MetaData{Message: msg},
	}
}
