package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type TcpProxy server
type UdpProxy server

type Config struct {
	TcpProxies  map[string]TcpProxy  `toml:"tcp"`
	UdpProxies  map[string]UdpProxy  `toml:"udp"`
	HttpProxies map[string]HttpProxy `toml:"http"`
}

type HttpProxy struct {
	Listen    string
	BasicAuth []BasicAuth `toml:"basic_auth"`
}

type server struct {
	From string
	To   string
}

type BasicAuth struct {
	User     string
	Password string
}

func main() {
	var config Config
	if _, err := toml.DecodeFile("./config.toml", &config); err != nil {
		log.Fatal("failed to parse config", err)
	}

	for _, proxy := range config.TcpProxies {
		go proxy.Start()
	}

	for _, proxy := range config.UdpProxies {
		go proxy.Start()
	}

	for _, proxy := range config.HttpProxies {
		log.Printf("http %+v", proxy)
		go proxy.Start()
	}

	select {}
}
