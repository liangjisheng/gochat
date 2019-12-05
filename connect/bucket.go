package connect

import (
	"gochat/proto"
	"sync"
)

// Bucket ...
type Bucket struct {
	cLock         sync.RWMutex     // protect the channels for chs
	chs           map[int]*Channel // map sub key to a channel
	bucketOptions BucketOptions
	rooms         map[int]*Room // bucket room channels
	routines      []chan *proto.PushRoomMsgRequest
	routinesNum   uint64
	broadcast     chan []byte
}

// BucketOptions ...
type BucketOptions struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}
