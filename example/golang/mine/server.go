package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"	
)


type Client struct {
	conn 	 net.Conn
	nickname string
	ch 		 chan string	
}

func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}
	
	
	msgchan := make(chan string, 1)
	addchan := make(chan Client, 1)
	rmchan  := make(chan Client, 1) // remove
	

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
		ch <- fmt.Sprintf("%s: %s", c.nickname, line)
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
	return
}


func promtpNick(c net.Conn, bufc *bufio.Reader) string{
	io.WriteString(c, "Welcome to the chat! What is your nick?\n")
	nick, _, _ := bufc.ReadLine()
	return string(nick)
}


func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- Client) {
	bufc := bufio.NewReader(c)
	defer c.Close()
	
	client := Client {
		conn:     c,
		nickname: promtpNick(c, bufc),
		ch:       make(chan string, 1),
	}
	
	if strings.TrimSpace(client.nickname) == "" {
		io.WriteString(c, "Invalid Username\n")
		return	
	}
	
	
	addchan <- client
	defer func(){
		msgchan <- fmt.Sprintf("User %s left the chat room.\n", client.nickname)
		log.Printf("Connection from %v closed.\n", c.RemoteAddr())
		rmchan <- client
	}()
	io.WriteString(c, fmt.Sprintf("Welcome, %s!\n\n", client.nickname))
	msgchan <- fmt.Sprintf("New user %s has joined the chat room.\n", client.nickname)

	// I/O	
	go client.WriteLinesFrom(client.ch)
	client.ReadLinesInto(msgchan)

}


func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan Client) {
	clients := make(map[net.Conn]chan<- string)
	
	for {
		select {
			case msg :=	<-msgchan:
				//log.Printf("New message: %s", msg)
				done := make(chan bool, len(clients))
				for _, ch := range clients {
					// go func(mch chan<- string) {
					// 	mch <- msg
					// } (ch)	
					// //ch <- msg
					go sendMsg(ch, msg, done)
				}

				fmt.Printf("Len = %d", len(clients))
				for i := 0; i < len(clients); i++ {
					fmt.Printf("i = %d", i)
					<- done
				}


			case client := <-addchan:
				fmt.Printf("New client: %v; nick: %s\n", client.conn, client.nickname)
				clients[client.conn] = client.ch
				
			case client := <-rmchan:
				log.Printf("Client disconnects: %v; nick: %s\n", client.conn, client.nickname)
				delete(clients, client.conn)
		}	
	}
	
}


func sendMsg(mch chan<- string, msg string, done chan<- bool) {
	mch <- msg
	done <- true
}