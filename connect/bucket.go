package connect

import (
	"gochat/proto"
	"sync"
	"sync/atomic"
)

// Bucket ...
type Bucket struct {
	// protect the channels for chs
	cLock sync.RWMutex
	// map sub key to a channel 存储所有的用户连接会话
	chs           map[int]*Channel
	bucketOptions BucketOptions
	// bucket room channels 存储所有的房间会话
	rooms map[int]*Room
	// goroutine 监听的所有 channel
	routines    []chan *proto.PushRoomMsgRequest
	routinesNum uint64
	broadcast   chan []byte
}

// BucketOptions ...
type BucketOptions struct {
	ChannelSize   int
	RoomSize      int
	RoutineAmount uint64
	RoutineSize   int
}

// NewBucket ...
func NewBucket(bucketOptions BucketOptions) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[int]*Channel, bucketOptions.ChannelSize)
	b.bucketOptions = bucketOptions
	b.routines = make([]chan *proto.PushRoomMsgRequest, bucketOptions.RoutineAmount)
	b.rooms = make(map[int]*Room, bucketOptions.RoomSize)
	for i := uint64(0); i < b.bucketOptions.RoutineAmount; i++ {
		c := make(chan *proto.PushRoomMsgRequest, bucketOptions.RoutineSize)
		b.routines[i] = c
		go b.PushRoom(c)
	}
	return
}

// PushRoom 向一个房间推送消息就是想一个房间中的所用用户推送消息
func (b *Bucket) PushRoom(ch chan *proto.PushRoomMsgRequest) {
	for {
		var (
			arg  *proto.PushRoomMsgRequest
			room *Room
		)
		arg = <-ch
		if room = b.Room(arg.RoomID); room != nil {
			room.Push(&arg.Msg)
		}
	}
}

// Room 通过房间id获取房间信息
func (b *Bucket) Room(rid int) (room *Room) {
	b.cLock.RLock()
	room, _ = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

// Put 通过 userID 和 roomID 把用户连接会话放入 bucket 中
func (b *Bucket) Put(userID int, roomID int, ch *Channel) (err error) {
	var (
		room *Room
		ok   bool
	)
	b.cLock.Lock()
	if roomID != NoRoom {
		if room, ok = b.rooms[roomID]; !ok {
			room = NewRoom(roomID)
			b.rooms[roomID] = room
		}
		ch.Room = room
	}
	ch.userID = userID
	b.chs[userID] = ch
	b.cLock.Unlock()

	if room != nil {
		err = room.Put(ch)
	}
	return
}

// DeleteChannel 从 bucket 中删除一个用户会话
func (b *Bucket) DeleteChannel(ch *Channel) {
	var (
		ok   bool
		room *Room
	)
	b.cLock.RLock()
	if ch, ok = b.chs[ch.userID]; ok {
		room = b.chs[ch.userID].Room
		// delete from bucket
		delete(b.chs, ch.userID)
	}
	if room != nil && room.DeleteChannel(ch) {
		// if room empty delete,will mark room.drop is true
		if room.drop == true {
			delete(b.rooms, room.ID)
		}
	}
	b.cLock.RUnlock()
}

// Channel ...
func (b *Bucket) Channel(userID int) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[userID]
	b.cLock.RUnlock()
	return
}

// BroadcastRoom ...
func (b *Bucket) BroadcastRoom(pushRoomMsgReq *proto.PushRoomMsgRequest) {
	num := atomic.AddUint64(&b.routinesNum, 1) % b.bucketOptions.RoutineAmount
	b.routines[num] <- pushRoomMsgReq
}
