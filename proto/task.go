package proto

// RedisMsg ...
type RedisMsg struct {
	Op           int               `json:"op"`
	ServerID     int               `json:"serverId,omitempty"`
	RoomID       int               `json:"roomId,omitempty"`
	UserID       int               `json:"userId,omitempty"`
	Msg          []byte            `json:"msg"`
	Count        int               `json:"count"`
	RoomUserInfo map[string]string `json:"roomUserInfo"`
}

// RedisRoomInfo ...
type RedisRoomInfo struct {
	Op           int               `json:"op"`
	RoomID       int               `json:"roomId,omitempty"`
	Count        int               `json:"count,omitempty"`
	RoomUserInfo map[string]string `json:"roomUserInfo"`
}

// RedisRoomCountMsg ...
type RedisRoomCountMsg struct {
	Count int `json:"count,omitempty"`
	Op    int `json:"op"`
}

// SuccessReply ...
type SuccessReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
