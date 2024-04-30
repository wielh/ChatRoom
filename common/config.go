package common

import (
	"github.com/go-pg/pg/v10"
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

func InitDB() {
	DB = pg.Connect(&pg.Options{
		User:     postgressUser,
		Password: postgresPassword,
		Database: postgresDBName,
	})

	if DB == nil {
		panic("db is nil")
	}
}
