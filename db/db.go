package db

import (
	"gochat/config"
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
		logrus.Error("connect db fail:%s", e.Error())
	}
}

func init() {
	initDB("gochat")
}
