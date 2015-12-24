package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"./rsa"
)

type Client struct {
	conn      net.Conn
	nickname  string
	ch        chan string
	RSAdecode rsa.RSA
	RSAencode rsa.RSA
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	msgchan := make(chan string, 1)
	addchan := make(chan Client, 1)
	rmchan := make(chan Client, 1) // remove chan

	go handleMessages(msgchan, addchan, rmchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn, msgchan, addchan, rmchan)
	}

}

func (c Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(c.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			return
		}
		ch <- line
	}
	return
}

func (c Client) WriteLinesFrom(ch <-chan string) {
	for {
		select {
		case msg := <-ch:
			_, err := io.WriteString(c.conn, msg)
			if err != nil {
				return
			}
		}
	}
}

func promtpNick(c net.Conn, bufc *bufio.Reader) string {
	io.WriteString(c, "Hi! Your nick:\n")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()

	client := Client{
		conn:     c,
		nickname: promtpNick(c, bufc),
		ch:       make(chan string, 1),
	}
	
	client.RSAencode = rsa.NewRSA()
	client.ch <- fmt.Printf("rsa %d %d", client.RSAencode.numb, client.RSAencode.encodeExp)
	

	if strings.TrimSpace(client.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return
	}

	addchan <- client
	defer func(){
		rmchan <- client
	}()

	// I/O
	go client.WriteLinesFrom(client.ch)
	client.ReadLinesInto(msgchan)
}

func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)

	for {
		select {
		case msg := <-msgchan:
			done := make(chan bool, len(clients))
			for _, ch := range clients {
				go sendMsg(ch, msg, done)
			}

			for i := 0; i < len(clients); i++ {
				<-done
			}

		case client := <-addchan:
			log.Printf("%s %v connects\n", client.nickname, client.conn.RemoteAddr())
			clients[client.conn] = client.ch

		case client := <-rmchan:
			log.Printf("%s %v disconnects\n", client.nickname, client.conn.RemoteAddr())
			delete(clients, client.conn)
		}
	}

}

func sendMsg(mch chan<- string, msg string, done chan<- bool) {
	mch <- msg
	done <- true
}
