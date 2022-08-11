package main

type factory struct {
	fileBroker *fsBroker
	client     *tcpClient
}

func (f *factory) getTcpClient() *tcpClient {
	if f.client == nil {
		f.client = &tcpClient{}
	}
	return f.client
}

func (f *factory) getFileBroker() *fsBroker {
	if f.fileBroker == nil {
		f.fileBroker = &fsBroker{path: DEFAULT_RECEIVE_FOLDER_PATH, contentChans: make(map[int]chan *delivery)}
	}
	return f.fileBroker
}
