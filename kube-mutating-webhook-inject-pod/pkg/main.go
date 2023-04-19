package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func main() {
	var runOption webHookSvrOptions

	// get command line parameters
	flag.IntVar(&runOption.port, "port", 8443, "Webhook server port.")
	flag.StringVar(&runOption.certFile, "tlsCertFile", "/etc/webhook/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&runOption.keyFile, "tlsKeyFile", "/etc/webhook/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")
	//flag.StringVar(&runOption.sidecarCfgFile,  "sidecarCfgFile", "/etc/webhook/config/sidecarconfig.yaml", "File containing the mutation configuration.")
	flag.StringVar(&runOption.sidecarCfgFile, "sidecarCfgFile", "D:\\Workspace\\Go\\src\\kube-mutating-webhook-inject-pod\\config.yaml", "File containing the mutation configuration.")
	flag.Parse()

	sidecarConfig, err := loadConfig(runOption.sidecarCfgFile)
	if err != nil {
		glog.Errorf("Failed to load configuration: %v", err)
		return
	}
	glog.Infof("[sidecarConfig:%v]", sidecarConfig)
	// 加载tls x509证书
	pair, err := tls.LoadX509KeyPair(runOption.certFile, runOption.keyFile)
	if err != nil {
		glog.Errorf("Failed to load key pair: %v", err)
		return
	}
	// 定义webhookhttp server, 并构造
	webhooksvr := &webhookServer{
		sidecarConfig: sidecarConfig,
		server: &http.Server{
			Addr:      fmt.Sprintf(":%v", runOption.port),
			TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
		},
	}
	// handler并关联path
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", webhooksvr.serveMutate)
	webhooksvr.server.Handler = mux

	// start webhook server in new grountine
	go func() {
		if err := webhooksvr.server.ListenAndServeTLS("", ""); err != nil {
			glog.Errorf("Failed to listen and serve webhook server: %v", err)
		}
	}()

	// listening OS shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	glog.Infof("Got 0S shutdown signal, shutting down webhook server gracefully...")
	webhooksvr.server.Shutdown(context.Background())
}
