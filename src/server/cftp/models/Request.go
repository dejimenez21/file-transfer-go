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

func NewDeliverRequest(senderAddress string, requestId int64, channel string, fileInfo File) *Request {
	return &Request{
		Method:   REQ_DELIVER,
		Meta:     MetaData{SenderAddress: senderAddress, RequestId: requestId},
		Channels: []string{channel},
		FileInfo: fileInfo,
	}
}
