package tools

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/bwmarrin/snowflake"
)

// SessionPrefix ...
const SessionPrefix = "sess_"

// GetSnowflakeID ...
func GetSnowflakeID() string {
	//default node id eq 1,this can modify to different serverId node
	node, _ := snowflake.NewNode(1)
	// Generate a snowflake ID.
	id := node.Generate().String()
	return id
}

// GetRandomToken ...
func GetRandomToken(length int) string {
	r := make([]byte, length)
	io.ReadFull(rand.Reader, r)
	return base64.URLEncoding.EncodeToString(r)
}

// CreateSessionID ...
func CreateSessionID(sessionID string) string {
	return SessionPrefix + sessionID
}

// GetSessionIDByUserID ...
func GetSessionIDByUserID(userID int) string {
	return fmt.Sprintf("sess_map_%d", userID)
}

// GetSessionName ...
func GetSessionName(sessionID string) string {
	return SessionPrefix + sessionID
}

// Sha1 ...
func Sha1(s string) (str string) {
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// GetNowDateTime ...
func GetNowDateTime() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}
