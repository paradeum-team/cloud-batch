package configs

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type cloudpods struct {
	BaseUrl              string `split_words:"true" required:"true"`
	Username             string `split_words:"true" required:"true"`
	Password             string `split_words:"true" required:"true"` // base64编译后密码
	Domain               string `split_words:"true" default:"Default"`
	VpcCidrBlock         string `split_words:"true" default:"172.30.0.0/16"`
	VpcPreName           string `split_words:"true" default:"cloudBatch"`
	NetworkGuestIpPrefix string `split_words:"true" default:"172.30"`
	NetworkPreName       string `split_words:"true" default:"cloudBatch"`
	AliyunDefaultImageId string `split_words:"true" default:"6851cdd4-95a4-4484-8b98-a61c6e61164a"`
	DefaultSecgroup      string `split_words:"true" default:"sg-cloudbatch-all"`
	DefaultKeypair       string `split_words:"true" default:"masterkey"`
}

var Cloudpods cloudpods

func init() {
	err := envconfig.Process("cloudpods", &Cloudpods)
	if err != nil {
		log.Fatalf("envconfig.Process cloudpods err: %+v", err)
	}
}
