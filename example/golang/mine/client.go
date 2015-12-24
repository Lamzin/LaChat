package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"./rsa"
)

type server struct {
	conn net.Conn
	encode rsa.RSA
	decode rsa.RSA
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		fmt.Print(err)
		return
	}
	handleConnection(conn)
}

func sendMessage(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}
		fmt.Fprintf(conn, text)
	}
}

func receiveMessage(conn net.Conn) {
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}
		
		arr := strings.Split(msg, " ")
		
		if len(arr) > 0 && arr[0] == "rsa" {
				
		} else if len(arr) > 0 && arr[0] == "msg" {
			fmt.Print(msg + "$ ")
		}
	}
}

func handleConnection(conn net.Conn) {
	go receiveMessage(conn)
	time.Sleep(time.Second * 1)
	sendMessage(conn)
}
