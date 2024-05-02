package common

import (
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

var DB *pg.DB
var ErrNoRows = pg.ErrNoRows
var postgressUser string = "postgres"
var postgresPassword string = "root"
var postgresDBName string = "chatroom"

var JWTSecretKey []byte = []byte("0d00-0721")

type UserType int

const (
	GoogleUser UserType = iota // 0
)

const (
	GatePort         = 80
	MicroAccountPort = 12300
	MicroRoomPort    = 12301
	MicroChatPort    = 12302
	PostgresPort     = 5432
	ElasticPort      = 9200
)

func initDB() *pg.DB {
	DB = pg.Connect(&pg.Options{
		User:       postgressUser,
		Password:   postgresPassword,
		Database:   postgresDBName,
		MaxRetries: 50,
	})

	return DB
}

func loggerInit() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func ConfigInit() {
	loggerInit()
	db := initDB()
	if db == nil {
		logrus.Error(generateMessage("common", "ConfigInit", "micro-account connect to postgre failed", nil))
		return
	}
}
