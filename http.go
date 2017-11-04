package main

import (
	"log"
	"net"
	"net/http"
	"time"
)

func (p HttpProxy) Start() {
	log.Printf("proxy is %v\n", p)
	srv := http.Server{
		Addr:         p.Listen,
		Handler:      http.HandlerFunc(httpHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("http failed to listen on %s\n", p.Listen)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "CONNECT" {
		w.Write([]byte("wrong request!"))
		return
	}

	target, err := net.DialTimeout("tcp", r.Host, 5*time.Second)
	if err != nil {
		log.Printf("failed to connect host%s\n", r.Host)
		return
	}
	defer target.Close()

	// 	w.WriteHeader(http.StatusOK)

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver don't support hijack", http.StatusInternalServerError)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	conn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))

	duplexCopy(conn, target)
}
