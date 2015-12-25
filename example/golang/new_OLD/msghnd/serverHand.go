package msghnd


import (
	"bufio"
	"os"
	"fmt"
	"../connection"
)


type ConnServer struct {
	connection.Connection
}


func (c ConnServer) HandlerIO() {
	fmt.Print("handlerIO\n")
	// go c.SendMsgServer()
	go c.ReceiveMsg()
	go c.handlerOut()
	c.handlerIn()
}


func (c ConnServer) handlerIn() {
	fmt.Print("handlerIn\n")
	reader := bufio.NewReader(os.Stdin)
	for c.Conn != nil {
		fmt.Print("handlerIn\n")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}
		c.Out <- text
	}
}


func (c ConnServer) handlerOut() {
	fmt.Print("handlerOut\n")
	for c.Conn != nil {
		fmt.Print("handlerOut\n")
		select {
		case msg := <- c.In:
			fmt.Print(msg + "$")
		}
	}
}
