package main

import "net"
import "fmt"
import "bufio"
//import "os"
import "time"
import "math/rand"

func main() {

	for i := 0; i < 300; i++ {

		go func(){
			// connect to this socket
			conn, err := net.Dial("tcp", "127.0.0.1:6000")
			
			if err != nil {
				fmt.Print("Here!")
				fmt.Print(err)	
				return
			}
				
			handleConnection(conn)
		}()
	}

	time.Sleep(time.Second * 120)

}



func sendMessage(conn net.Conn) {
	//reader := bufio.NewReader(os.Stdin)

	for {
		//fmt.Print("Text to send: ")	

		text := fmt.Sprintf("test bot #%d\n", rand.Int31())
		time.Sleep(time.Second * 1)
	
//		text, err := reader.ReadString('\n')
//		if err != nil {
//			fmt.Print(err)
//			return
//		}
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
		
		fmt.Print(msg)
	}
}


func handleConnection(conn net.Conn) {
	go receiveMessage(conn)
	time.Sleep(time.Second * 1)
	sendMessage(conn)
}