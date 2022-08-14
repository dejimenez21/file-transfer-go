package main

type factory struct {
	fileBroker *fsBroker
	client     *tcpClient
	serverAddr string
}

func (f *factory) getTcpClient() *tcpClient {
	if f.client == nil {
		f.client = &tcpClient{serverAddr: f.serverAddr}
	}
	return f.client
}

func (f *factory) getFileBroker() *fsBroker {
	if f.fileBroker == nil {
		f.fileBroker = &fsBroker{path: DEFAULT_RECEIVE_FOLDER_PATH, contentChans: make(map[int]chan *delivery)}
	}
	return f.fileBroker
}
