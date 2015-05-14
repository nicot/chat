package main

import (
	"bytes"
	"log"
	"net"
	"strings"
)

type message struct {
	user string
	m    string
}

type sub struct {
	user string
	out  chan message
}

func clean(b []byte) string {
	b = bytes.Trim(b, "\x00")
	s := strings.TrimSpace(string(b))
	return s
}

func read(conn net.Conn, user string, in chan<- message) {
	// Error handling and EOF
	for {
		b := make([]byte, 256) // TODO learn to zero a buffer
		_, e := conn.Read(b)
		if e != nil {
			log.Println("ERROR: ", e)
			break
		}
		in <- message{clean(b), user}
	}
}

func subscribe(s chan sub, in chan message) {
	subs := make([]sub, 0)
	for {
		select {
		case subscriber := <-s:
			go func() {
				in <- message{"server", subscriber.user + " logged in"}
				var us string
				for _, u := range subs {
					us += "\n" + u.user
				}
				subscriber.out <- message{"server", "Online:" + us}
				subs = append(subs, subscriber)
			}()
		case m := <-in:
			for _, subscriber := range subs {
				subscriber.out <- m
			}
		}
	}
}

func handle(conn net.Conn, in chan<- message, out chan message, subs chan sub) {
	w := "Welcome to chat!\n"
	conn.Write([]byte(w))
	w = "What do you want your nickname to be?\n"
	conn.Write([]byte(w))
	b := make([]byte, 16)
	conn.Read(b)
	user := clean(b)
	subs <- sub{user, out}
	go read(conn, user, in)
	for m := range out {
		conn.Write([]byte(m.user + ": " + m.m + "\n"))
	}
}

func main() {
	port := ":3030"
	proto := "tcp"
	ln, err := net.Listen(proto, port)
	if err != nil {
		log.Fatal("Couldn't listen on %s:\n%s", port, err)
	}
	in := make(chan message)
	s := make(chan sub)
	go subscribe(s, in)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		out := make(chan message)
		go handle(conn, in, out, s)
	}
}
