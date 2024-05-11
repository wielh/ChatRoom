package common

import (
	"net"
	"net/http"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var DB *pg.DB
var ErrNoRows = pg.ErrNoRows
var ElasticClient *elasticsearch.Client
var postgressUser string = "postgres"
var postgresPassword string = "root"
var postgresDBName string = "chatroom"

var JWTSecretKey []byte = []byte("0d00-0721")

type UserType int

const (
	None       UserType = -1
	GoogleUser UserType = 0
)

const (
	GatePort         = 80
	MicroAccountPort = 12300
	MicroRoomPort    = 12301
	MicroChatPort    = 12302
	PostgresPort     = 5432
	ElasticPort      = 9200
)

var GoogleOauth2Config = &oauth2.Config{
	ClientID:     "***.apps.googleusercontent.com",
	ClientSecret: "***",
	RedirectURL:  "http://127.0.0.1:80/account/google_callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	},
	Endpoint: google.Endpoint,
}

func DBInit() *pg.DB {
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

func elasticInit() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200",
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		},
	}

	return elasticsearch.NewClient(cfg)
}

func ConfigInit() {
	loggerInit()

	var err error
	ElasticClient, err = elasticInit()
	if err != nil {
		logrus.Error(generateMessage("common", "ConfigInit", "connect to elastic failed", nil))
		return
	}

	db := DBInit()
	if db == nil {
		logrus.Error(generateMessage("common", "ConfigInit", "connect to postgre failed", nil))
		return
	}
}
