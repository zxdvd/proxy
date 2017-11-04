package main

import (
	"encoding/base64"
	"io"
	"log"
	"strings"
	"sync"
)

func ioCopy(wg *sync.WaitGroup, dst io.Writer, src io.Reader) {
	defer wg.Done()
	if _, err := io.Copy(dst, src); err != nil {
		log.Printf("tcp: ioCopy err, %v\n", err)
	}
}

type Conn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

func duplexCopy(conn1, conn2 Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go ioCopy(&wg, conn1, conn2)
	go ioCopy(&wg, conn2, conn1)
	wg.Wait()
}

// copy from net/http/request.go
func parseBasicAuth(auth string) (user, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
