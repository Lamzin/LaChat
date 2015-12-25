package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	// "strings"	
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
	return
}

func handleConnection(c net.Conn, 
						msgchan chan<- string, 
						addchan chan<- Client, 
						rmchan chan<- Client) {
	defer c.Close()
	
	client := Client {
		conn:     c,
		nickname: "Oleh",//promtpNick(c, bufc),
		ch:       make(chan string, 1),
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
			case msg :=	<-msgchan:
				log.Printf("New message: %s", msg)
				for _, ch := range clients {
					ch <- msg
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
