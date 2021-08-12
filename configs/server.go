package configs

import (
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type server struct {
	ListenAddr          string `split_words:"true" default:"0.0.0.0:5140"`
	DefaultReadTimeOut  int    `split_words:"true" default:"60"`
	DefaultWriteTimeOut int    `split_words:"true" default:"60"`
	IsOpenMetrics       bool   `split_words:"true" default:"true"`
	MetricsAddr         string `split_words:"true" default:"0.0.0.0:20001"`
	RuntimeRootPath     string `split_words:"true" default:"runtime"`
	MediaSavePath       string `split_words:"true" default:"upload/files/"`
	LevelDBPath         string `split_words:"true" default:"leveldb/db"`
	JwtSecret           string `split_words:"true" required:"true"`
	HostName            string `ignored:"true"`
	Version             string `ignored:"true"`
}

type logConfig struct {
	IsOutPutFile bool   `split_words:"true" default:"true"`
	Level        string `split_words:"true" default:"info"`
	Path         string `split_words:"true" default:"logs"`
	MaxAgeDay    int    `split_words:"true" default:"7"`
	RotationDay  int    `split_words:"true" default:"1"`
	TimeFormat   string `split_words:"true" default:"2006-01-02 15:04:05"`
	AppName      string `split_words:"true" default:"cloudBatch"`
}

var LogConfig logConfig
var Server server

func init() {
	err := envconfig.Process("server", &Server)
	if err != nil {
		log.Fatalf("envconfig.Process server err: %+v", err)
	}

	err = envconfig.Process("log", &LogConfig)
	if err != nil {
		log.Fatalf("envconfig.Process log err: %+v", err)
	}

	Server.HostName, err = os.Hostname()
	if err != nil {
		log.Fatalf("Server.HostName get err: %+v", err)
	}

	dat, err := ioutil.ReadFile("./VERSION")
	if err != nil {
		log.Printf("read VERSION failed: %v", err)
	} else {
		Server.Version = strings.TrimSpace(string(dat))
		log.Printf("VERSION: %s", Server.Version)
	}
}
