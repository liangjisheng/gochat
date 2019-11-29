package proto

// Msg ...
type Msg struct {
	Ver       int    `json:"ver"`  // protocol version
	Operation int    `json:"op"`   // operation for request
	SeqID     string `json:"seq"`  // sequence number chosen by client
	Body      []byte `json:"body"` // binary body bytes
}

// PushMsgRequest ...
type PushMsgRequest struct {
	UserID int
	Msg    Msg
}

// PushRoomMsgRequest ...
type PushRoomMsgRequest struct {
	RoomID int
	Msg    Msg
}

// PushRoomCountRequest ...
type PushRoomCountRequest struct {
	RoomID int
	Count  int
}
