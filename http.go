package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func (auth BasicAuth) validate(user, password string) bool {
	if auth.Password == password {
		if auth.User == "" || auth.User == user {
			return true
		}
	}
	return false
}

func (p HttpProxy) Start() {
	log.Printf("http proxy: %+v\n", p)
	srv := http.Server{
		Addr:         p.Listen,
		Handler:      http.HandlerFunc(p.httpHandler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("http failed to listen on %s\n", p.Listen)
	}
	if p.certFile != "" && p.keyFile != "" {
		err = srv.ListenAndServeTLS(p.certFile, p.keyFile)
		if err != nil {
			log.Printf("https failed on %s\n", p.Listen)
			log.Printf("https error: %v\n", err)
		}
	}
}

func (p HttpProxy) shouldCheckBasicAuth() bool {
	return len(p.BasicAuth) > 0
}

func (p HttpProxy) checkBasicAuth(r *http.Request) bool {
	if !p.shouldCheckBasicAuth() {
		return true
	}

	auth := r.Header.Get("Proxy-Authorization")

	user, password, ok := parseBasicAuth(auth)
	if !ok {
		return ok
	}
	for _, auth_user := range p.BasicAuth {
		if auth_user.validate(user, password) {
			return true
		}
	}
	return false
}

func empty(req *http.Request) {}

func HttpForwardProxy(w http.ResponseWriter, r *http.Request) {
	p := &httputil.ReverseProxy{
		Director: empty,
	}
	p.ServeHTTP(w, r)
}

func (p HttpProxy) httpHandler(w http.ResponseWriter, r *http.Request) {
	if !p.checkBasicAuth(r) {
		log.Printf("Authentication failed")
		http.Error(w, "auth failed!", 400)
		return
	}

	if r.Method != "CONNECT" {
		HttpForwardProxy(w, r)
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
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

	duplexCopy(conn, target)
}
