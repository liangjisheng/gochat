package connect

import (
	"errors"
	"gochat/proto"
	"sync"
)

// NoRoom ...
const NoRoom = -1

// Room ...
type Room struct {
	ID          int
	OnlineCount int // room online user count
	rLock       sync.RWMutex
	drop        bool // make room is live
	next        *Channel
}

// NewRoom ...
func NewRoom(roomID int) *Room {
	room := new(Room)
	room.ID = roomID
	room.drop = false
	room.next = nil
	room.OnlineCount = 0
	return room
}

// Put 某个用户进入房间
func (r *Room) Put(ch *Channel) (err error) {
	// doubly linked list
	r.rLock.Lock()
	defer r.rLock.Unlock()
	if !r.drop {
		if r.next != nil {
			r.next.Prev = ch
		}
		ch.Next = r.next
		ch.Prev = nil
		r.next = ch
		r.OnlineCount++
	} else {
		err = errors.New("room drop")
	}
	return
}

// Push 向一个房间中的所有用户推送消息
func (r *Room) Push(msg *proto.Msg) {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Push(msg)
	}
	r.rLock.RUnlock()
	return
}

// DeleteChannel 某个用户离开房间
func (r *Room) DeleteChannel(ch *Channel) bool {
	r.rLock.RLock()
	if ch.Next != nil {
		//if not footer
		ch.Next.Prev = ch.Prev
	}
	if ch.Prev != nil {
		// if not header
		ch.Prev.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	r.OnlineCount--
	r.drop = false
	if r.OnlineCount <= 0 {
		r.drop = true
	}
	r.rLock.RUnlock()
	return r.drop
}
