package connection


import (
	"fmt"
	"bufio"
	"net"
	// "io"
)


type Connection struct {
	Conn net.Conn
	In   chan string
	Out  chan string
}


func (c Connection) Connect(addres string) { //"127.0.0.1:6000"
	if conn, err := net.Dial("tcp", addres); err != nil {
		fmt.Print(err)
	} else {
		c.Conn = conn
		fmt.Print("Ok!!!\n")
	}
}


func (c Connection) SendMsgClient() {
	for c.Conn != nil {
		select {
		case msg := <- c.Out:
			fmt.Fprintf(c.Conn, msg)
		}
	}
}


func (c Connection) SendMsgServer() {
	for c.Conn != nil {
		select {
		case msg := <- c.Out:
			c.Conn.Write([]byte(msg + "\n"))
		}
	}
}


func (c Connection) ReceiveMsg() {
	reader := bufio.NewReader(c.Conn)
	for c.Conn != nil {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}
		c.In <- msg
	}
}
