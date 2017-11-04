package main

import (
	"log"
	"net"
	"sync"
)

func (p UdpProxy) Start() {
	var fromAddr *net.UDPAddr
	var toAddr *net.UDPAddr
	fromAddr, err := net.ResolveUDPAddr("udp", p.From)
	if err != nil {
		log.Printf("wrong udp proxy address %s\n", p.From)
		return
	}
	toAddr, err = net.ResolveUDPAddr("udp", p.To)
	if err != nil {
		log.Printf("wrong udp upstream address %s\n", p.To)
		return
	}
	socket, err := net.ListenUDP("udp", fromAddr)
	if err != nil {
		log.Printf("failed to listen on %s. %v", p.From, err)
		return
	}
	defer socket.Close()
	log.Printf("udp: listen @ port %s\n", p.From)

	var clients sync.Map

	buf := make([]byte, 500)
	for {
		var client *net.UDPAddr
		n, client, err := socket.ReadFromUDP(buf)
		if err != nil {
			log.Printf("failed to read udp", err)
			continue
		}
		if len(buf) == 0 {
			continue
		}
		upstream, ok := clients.Load(client.String())
		if !ok {
			upstream, err = net.DialUDP("udp", nil, toAddr)
			if err != nil {
				log.Printf("failed to dial upstream", err)
				continue
			}
			clients.Store(client.String(), upstream)
		}
		upstream.(*net.UDPConn).Write(buf[:n])
		go RecvUpstream(upstream.(*net.UDPConn), socket, client)
	}
}

func RecvUpstream(upstream *net.UDPConn, frontend *net.UDPConn, client *net.UDPAddr) {
	buf := make([]byte, 500)
	for {
		n, err := upstream.Read(buf)
		if err != nil {
			log.Printf("err read", err)
			return
		}
		if len(buf) == 0 {
			continue
		}
		if _, err := frontend.WriteToUDP(buf[:n], client); err != nil {
			log.Printf("write failed", err)
			return
		}
	}
}
