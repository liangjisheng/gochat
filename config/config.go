package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var once sync.Once
var realPath string

// Conf ...
var Conf *Config

const (
	SuccessReplyCode      = 0
	FailReplyCode         = 1
	SuccessReplyMsg       = "success"
	QueueName             = "gochat_sub"
	RedisBaseValidTime    = 86400
	RedisPrefix           = "gochat_"
	RedisRoomPrefix       = "gochat_room_"
	RedisRoomOnlinePrefix = "gochat_room_online_count_"
	MsgVersion            = 1
	OpSingleSend          = 2 // single user
	OpRoomSend            = 3 // send to room
	OpRoomCountSend       = 4 // get online user count
	OpRoomInfoSend        = 5 // send info to room
)

// Config ...
type Config struct {
	Common  Common
	Connect ConnectConfig
	Logic   LogicConfig
	Task    TaskConfig
	API     APIConfig
	Site    SiteConfig
}

func init() {
	Init()
}

// Init ...
func Init() {
	once.Do(func() {
		env := GetMode()
		realPath, _ := filepath.Abs("./")
		configFilePath := realPath + "/config/" + env + "/"
		viper.SetConfigType("toml")
		viper.AddConfigPath(configFilePath)

		viper.SetConfigName("/connect")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/common")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/task")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/logic")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/api")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("/site")
		err = viper.MergeInConfig()
		if err != nil {
			panic(err)
		}
		Conf = new(Config)
		viper.Unmarshal(&Conf.Common)
		viper.Unmarshal(&Conf.Connect)
		viper.Unmarshal(&Conf.Task)
		viper.Unmarshal(&Conf.Logic)
		viper.Unmarshal(&Conf.API)
		viper.Unmarshal(&Conf.Site)

		// fmt.Printf("config common: %+v\n", Conf.Common)
		// fmt.Printf("config connect: %+v\n", Conf.Connect)
		// fmt.Printf("config task: %+v\n", Conf.Task)
		// fmt.Printf("config login: %+v\n", Conf.Logic)
		// fmt.Printf("config api: %+v\n", Conf.API)
		// fmt.Printf("config site: %+v\n", Conf.Site)
	})
}

// GetMode ...
func GetMode() string {
	env := os.Getenv("RUN_MODE")
	if env == "" {
		env = "dev"
	}
	return env
}

// GetGinRunMode ...
func GetGinRunMode() string {
	env := GetMode()
	// gin have debug,test,release mode
	if env == "dev" {
		return "debug"
	}
	if env == "test" {
		return "debug"
	}
	if env == "prod" {
		return "release"
	}
	return "release"
}

// CommonEtcd ...
type CommonEtcd struct {
	Host              string `mapstructure:"host"`
	BasePath          string `mapstructure:"basePath"`
	ServerPathLogic   string `mapstructure:"serverPathLogic"`
	ServerPathConnect string `mapstructure:"serverPathConnect"`
	ServerID          int    `mapstructure:"serverId"`
}

// CommonRedis ...
type CommonRedis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}

// Common ...
type Common struct {
	CommonEtcd  CommonEtcd  `mapstructure:"common-etcd"`
	CommonRedis CommonRedis `mapstructure:"common-redis"`
}

// ConnectBase ...
type ConnectBase struct {
	ServerID int    `mapstructure:"serverId"`
	CertPath string `mapstructure:"certPath"`
	KeyPath  string `mapstructure:"keyPath"`
}

// ConnectWebsocket ...
type ConnectWebsocket struct {
	Bind string `mapstructure:"bind"`
}

// ConnectRPCAddress ...
type ConnectRPCAddress struct {
	Address string `mapstructure:"address"`
}

// ConnectBucket ...
type ConnectBucket struct {
	CPUNum        int    `mapstructure:"cpuNum"`
	Channel       int    `mapstructure:"channel"`
	Room          int    `mapstructure:"room"`
	SrvProto      int    `mapstructure:"svrProto"`
	RoutineAmount uint64 `mapstructure:"routineAmount"`
	RoutineSize   int    `mapstructure:"routineSize"`
}

// ConnectConfig ...
type ConnectConfig struct {
	ConnectBase       ConnectBase       `mapstructure:"connect-base"`
	ConnectRPCAddress ConnectRPCAddress `mapstructure:"connect-rpcAddress"`
	ConnectBucket     ConnectBucket     `mapstructure:"connect-bucket"`
	ConnectWebsocket  ConnectWebsocket  `mapstructure:"connect-websocket"`
}

// LogicBase ...
type LogicBase struct {
	CPUNum     int    `mapstructure:"cpuNum"`
	RPCAddress string `mapstructure:"rpcAddress"`
	CertPath   string `mapstructure:"certPath"`
	KeyPath    string `mapstructure:"keyPath"`
}

// LogicRedis ...
// type LogicRedis struct {
// 	RedisAddress  string `mapstructure:"redisAddress"`
// 	RedisPassword string `mapstructure:"redisPassword"`
// }

// LogicEtcd ...
// type LogicEtcd struct {
// 	Host     string `mapstructure:"host"`
// 	BasePath string `mapstructure:"basePath"`
// 	ServerID string `mapstructure:"serverId"`
// }

// LogicConfig ...
type LogicConfig struct {
	LogicBase LogicBase `mapstructure:"logic-base"`
	// LogicRedis LogicRedis `mapstructure:"logic-redis"`
	// LogicEtcd  LogicEtcd  `mapstructure:"logic-etcd"`
}

// TaskBase ...
type TaskBase struct {
	CPUNum        int    `mapstructure:"cpuNum"`
	RedisAddr     string `mapstructure:"redisAddr"`
	RedisPassword string `mapstructure:"redisPassword"`
	RPCAddress    string `mapstructure:"rpcAddress"`
	PushChan      int    `mapstructure:"pushChan"`
	PushChanSize  int    `mapstructure:"pushChanSize"`
}

// TaskConfig ...
type TaskConfig struct {
	TaskBase TaskBase `mapstructure:"task-base"`
}

// APIBase ...
type APIBase struct {
	ListenPort int `mapstructure:"listenPort"`
}

// APIConfig ...
type APIConfig struct {
	APIBase APIBase `mapstructure:"api-base"`
}

// SiteBase ...
type SiteBase struct {
	ListenPort int `mapstructure:"listenPort"`
}

// SiteConfig ...
type SiteConfig struct {
	SiteBase SiteBase `mapstructure:"site-base"`
}
