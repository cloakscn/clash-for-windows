package main

import (
	"os"
	"os/signal"
	"syscall"

	C "github.com/cloakscn/clash-for-windows/constant"
	"github.com/cloakscn/clash-for-windows/hub"
	"github.com/cloakscn/clash-for-windows/proxy/http"
	"github.com/cloakscn/clash-for-windows/proxy/socks"
	"github.com/cloakscn/clash-for-windows/tunnel"

	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := C.GetConfig()
	if err != nil {
		log.Fatalf("Read config error: %s", err.Error())
	}

	port, socksPort := C.DefalutHTTPPort, C.DefalutSOCKSPort
	section := cfg.Section("General")
	if key, err := section.GetKey("port"); err == nil {
		port = key.Value()
	}

	if key, err := section.GetKey("socks-port"); err == nil {
		socksPort = key.Value()
	}

	err = tunnel.GetInstance().UpdateConfig()
	if err != nil {
		log.Fatalf("Parse config error: %s", err.Error())
	}

	go http.NewHttpProxy(port)
	go socks.NewSocksProxy(socksPort)

	// Hub
	if key, err := section.GetKey("external-controller"); err == nil {
		go hub.NewHub(key.Value())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
