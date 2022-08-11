package main

type delivery struct {
	ID      int64
	Seq     int64
	Size    int
	Content []byte
}

var nextDeliveryID int64 = 0
