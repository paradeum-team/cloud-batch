package configs

import (
	"github.com/kataras/golog"
	"github.com/kelseyhightower/envconfig"
)

type redis struct {
	Host string `split_words:"true" required:"true"`
	Port int    `split_words:"true" default:"6379"`
	Pass string `split_words:"true" required:"false"`
	DB   int    `split_words:"true" required:"true"`
}

var Redis redis

func init() {
	err := envconfig.Process("redis", &Redis)
	if err != nil {
		golog.Fatalf("envconfig.Process redis err: %v", err)
	}
}
