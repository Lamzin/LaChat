package msghnd


import (
	"bufio"
	"os"
	"fmt"
	"../connection"
)


type ConnClient struct {
	connection.Connection
}


func (c ConnClient) HandlerIO() {
	go c.SendMsgClient()
	go c.ReceiveMsg()
	go c.handlerIn()
	c.handlerOut()
}


func (c ConnClient) handlerIn() {
	reader := bufio.NewReader(os.Stdin)
	for c.Conn != nil {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}
		c.Out <- text + "\n"
	}
}


func (c ConnClient) handlerOut() {
	for c.Conn != nil {
		select {
		case msg := <- c.In:
			fmt.Print(msg + "$")
		}
	}
}
