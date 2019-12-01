package dao

import (
	"errors"
	"gochat/db"
	"time"
)

var dbIns = db.GetDb("gochat")

// User ...
type User struct {
	ID         int `gorm:"primary_key"`
	UserName   string
	Password   string
	CreateTime time.Time
	db.GoChat
}

// TableName ...
func (u *User) TableName() string {
	return "user"
}

// Add ...
func (u *User) Add() (userID int, err error) {
	if u.UserName == "" || u.Password == "" {
		return 0, errors.New("user_name or password empty")
	}
	oUser := u.CheckHaveUserName(u.UserName)
	if oUser.ID > 0 {
		return oUser.ID, nil
	}
	u.CreateTime = time.Now()
	if err = dbIns.Table(u.TableName()).Create(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}

// CheckHaveUserName ...
func (u *User) CheckHaveUserName(userName string) (data User) {
	dbIns.Table(u.TableName()).Where("user_name=?", userName).First(&data)
	return
}

// GetUserNameByUserID ...
func (u *User) GetUserNameByUserID(userID int) (userName string) {
	var data User
	dbIns.Table(u.TableName()).Where("user_id=?", userID).First(&data)
	return data.UserName
}
