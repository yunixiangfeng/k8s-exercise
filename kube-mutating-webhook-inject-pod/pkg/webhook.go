package main

import (
	"crypto/sha256"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
)

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
	Volumes    []corev1.Volume    `yaml:"volumes"`
}

// Webhook Server options
type webHookSvrOptions struct {
	port           int    // 监听https的端口
	certFile       string // https x509 证书路径
	keyFile        string // https x509 证书私钥路径
	sidecarCfgFile string // 注入sidecar容器的配置文件路径
}

type webhookServer struct {
	sidecarConfig *Config      // 注入sidecar容器的配置
	server        *http.Server // http serer
}

// webhookServer的mutate handler
func (ws *webhookServer) serveMutate(w http.ResponseWriter, r *http.Request) {

}

func loadConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	glog.Infof("New configuration: sha256sum %x", sha256.Sum256(data))
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
