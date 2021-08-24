package configs

import (
	"github.com/kelseyhightower/envconfig"
	"log"
)

type k8sbfs struct {
	// 安装 k8s 的 ansible 目录
	AnsibleConfigPath string `split_words:"true" default:"/etc/ansible"`
	AnsibleHostsName  string `split_words:"true" default:"hosts"`
	K8sBaseDir        string `split_words:"true" default:"/data/offline-k8s-install-package/offline-k8s"`
	K8sAddNodeScript  string `split_words:"true" default:"add_node.sh"`
	BfsPath           string `split_words:"true" default:"/root/bfs"`
	BfsValuesName     string `split_words:"true" default:"values.yaml"`
}

var K8sBfs k8sbfs

func init() {
	err := envconfig.Process("k8sbfs", &K8sBfs)
	if err != nil {
		log.Fatalf("envconfig.Process k8sbfs err: %+v", err)
	}
}
