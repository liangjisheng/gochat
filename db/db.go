package db

import (
	"fmt"
	"gochat/config"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	// sqlite ...
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbMap = map[string]*gorm.DB{}
var syncLock sync.Mutex

// GoChat ...
type GoChat struct {
}

// GetDbName ...
func (*GoChat) GetDbName() string {
	return "gochat"
}

// GetDb ...
func GetDb(dbName string) (db *gorm.DB) {
	if db, ok := dbMap[dbName]; ok {
		return db
	}
	return nil
}

// User ...
type User struct {
	ID    uint64 `gorm:"id"`
	Name  string `gorm:"user_name"`
	Pwd   string `gorm:"password"`
	CTime uint64 `gorm:"create_time"`
}

func initDB(dbName string) {
	var e error
	// if prod env , you should change mysql driver for yourself
	realPath, _ := filepath.Abs("./")
	configFilePath := realPath + "/db/gochat.sqlite3"
	syncLock.Lock()
	dbMap[dbName], e = gorm.Open("sqlite3", configFilePath)
	dbMap[dbName].DB().SetMaxIdleConns(4)
	dbMap[dbName].DB().SetMaxOpenConns(20)
	dbMap[dbName].DB().SetConnMaxLifetime(8 * time.Second)
	if config.GetMode() == "dev" {
		dbMap[dbName].LogMode(true)
	}
	syncLock.Unlock()
	if e != nil {
		logrus.Errorf("connect db fail:%s", e.Error())
	}
}

func queryAll() {
	gochat := &GoChat{}
	gochatDB := GetDb(gochat.GetDbName())
	sqlStr := fmt.Sprintf("select * from user")
	gochatDB = gochatDB.Exec(sqlStr)
	if gochatDB.Error != nil {
		log.Println("select db err:", gochatDB.Error)
		return
	}

	// var id, timestamp uint64
	// var username, userpwd string
	var users []User
	gochatDB = gochatDB.Table("user").Find(&users)
	if gochatDB.Error != nil {
		log.Println("find db err:", gochatDB.Error)
		return
	}

	for _, user := range users {
		fmt.Printf("user: %+v\n", user)
	}
}

func countDB() {
	gochat := &GoChat{}
	gochatDB := GetDb(gochat.GetDbName())
	sqlStr := "select count(*) from user"
	gochatDB = gochatDB.Exec(sqlStr)
	if gochatDB.Error != nil {
		log.Println("select db err:", gochatDB.Error)
		return
	}

	var count uint64
	gochatDB = gochatDB.Table("user").Count(&count)
	if gochatDB.Error != nil {
		log.Println("scan db err:", gochatDB.Error)
		return
	}
	fmt.Println("count:", count)
}

func init() {
	initDB("gochat")
	// queryAll()
	countDB()
}
