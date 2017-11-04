package main

import (
	"log"
	"net"
	"time"
)

func (p TcpProxy) Start() {
	socket, err := net.Listen("tcp", p.From)
	if err != nil {
		log.Printf("failed to listen on %s. %v\n", p.From, err)
		return
	}
	defer socket.Close()
	log.Printf("tcp: listen @ port %s\n", p.From)

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("failed to accept %s. %v", p.From, err)
		}
		go p.handleTcpConnection(conn)
	}
}

func (p TcpProxy) handleTcpConnection(conn net.Conn) {
	defer conn.Close()

	backend, err := net.DialTimeout("tcp", p.To, 10*time.Second)
	if err != nil {
		log.Printf("failed to connect to %s. %v\n", p.To, err)
		return
	}
	defer backend.Close()

	duplexCopy(conn, backend)
}
