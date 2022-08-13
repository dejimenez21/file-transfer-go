package models

type Delivery struct {
	ID      int64
	Seq     int64
	Size    int
	Content []byte
}

var nextDeliveryID int64 = 0

//TODO: Add constructor func

func NewDeliveryId() int64 {
	nextDeliveryID++
	return nextDeliveryID
}
