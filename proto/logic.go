package proto

// LoginRequest ...
type LoginRequest struct {
	Name     string
	Password string
}

// LoginResponse ...
type LoginResponse struct {
	Code      int
	AuthToken string
}

// GetUserInfoRequest ...
type GetUserInfoRequest struct {
	UserID int
}

// GetUserInfoResponse ...
type GetUserInfoResponse struct {
	Code     int
	UserID   int
	UserName string
}

// RegisterRequest ...
type RegisterRequest struct {
	Name     string
	Password string
}

// RegisterReply ...
type RegisterReply struct {
	Code      int
	AuthToken string
}

// LogoutRequest ...
type LogoutRequest struct {
	AuthToken string
}

// LogoutResponse ...
type LogoutResponse struct {
	Code int
}

// CheckAuthRequest ...
type CheckAuthRequest struct {
	AuthToken string
}

// CheckAuthResponse ...
type CheckAuthResponse struct {
	Code     int
	UserID   int
	UserName string
}

// ConnectRequest ...
type ConnectRequest struct {
	AuthToken string `json:"authToken"`
	RoomID    int    `json:"roomId"`
	ServerID  int    `json:"serverId"`
}

// ConnectReply ...
type ConnectReply struct {
	UserID int
}

// DisConnectRequest ...
type DisConnectRequest struct {
	RoomID int
	UserID int
}

// DisConnectReply ...
type DisConnectReply struct {
	Has bool
}

// Send ...
type Send struct {
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	FromUserID   int    `json:"fromUserId"`
	FromUserName string `json:"fromUserName"`
	ToUserID     int    `json:"toUserId"`
	ToUserName   string `json:"toUserName"`
	RoomID       int    `json:"roomId"`
	Op           int    `json:"op"`
	CreateTime   string `json:"createTime"`
}
