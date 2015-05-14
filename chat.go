package main

import (
	"log"
	"net"
)

type robustConn struct {
	conn net.Conn
}

func (c robustConn) write(b []byte) {
	t := 0
	for t < len(b) {
		// timeout
		n, err := c.conn.Write(b)
		t += n
		if err != nil {
			log.Println(err)
		}
	}
}

func (c robustConn) read() []byte {
	b := make([]byte, 1e5)
	//timeout
	n, err := c.conn.Read(b)
	if err != nil {
		log.Println(err)
	}
	r := make([]byte, n)
	copy(r, b)
	return r
}

func handle(c robustConn) {
	w := "Welcome to chat\n"
	c.write([]byte(w))
	for {
		b := c.read()
		log.Println(string(b))
	}
}

func main() {
	port := ":3030"
	proto := "tcp"
	ln, err := net.Listen(proto, port)
	if err != nil {
		log.Fatal("Couldn't listen on %s:\n%s", port, err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		c := robustConn{conn}
		go handle(c)
	}
}
