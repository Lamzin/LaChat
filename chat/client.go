package main

import "net"
import "fmt"
import "bufio"
import "os"
import "time"
import "strings"
import "strconv"
import "./rsa"

func main() {
	var dial dialog
	dial.is_exist = false
	dial.RSAdecode = rsa.NewRSA()
	dial.friend1, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	dial.friend2, _ = bufio.NewReader(os.Stdin).ReadString('\n')

	dial.friend1 = strings.TrimSpace(dial.friend1)
	dial.friend2 = strings.TrimSpace(dial.friend2)

	conn, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		fmt.Print("Here!")
		fmt.Print(err)	
		return
	}
	
	dial.conn = conn
	dial.handleConnection()

}


type dialog struct {
	conn net.Conn
	friend1 string
	friend2 string
	RSAdecode rsa.RSA
	RSAencode rsa.RSA
	keyDecode int
	keyEncode int
	is_exist bool
}


func (d dialog) sendMessage() {
	fmt.Fprintf(d.conn, "connect %s %s %d %d\n", 
				d.friend1, 
				d.friend2, 
				d.RSAdecode.Numb, 
				d.RSAdecode.DecodeExp)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}

		// if !d.is_exist {
		// 	fmt.Print("Wait your friend...")
		// 	continue
		// }

		key := 2
		fmt.Fprintf(d.conn, "text %d %s\n", d.RSAdecode.Decode(int64(key)), DecodeText(text, key))
	}
}


func (d dialog) receiveMessage() {
	for {
		msg, err := bufio.NewReader(d.conn).ReadString('\n')
		if err != nil {
			fmt.Print(err)
			return
		}			
		
		arr := strings.Split(msg, " ")
		if len(arr) > 2 {
			if arr[0] == "decode" {
				if len(arr) == 3 {
					if n, err1 := strconv.Atoi(arr[1]); err1 == nil {
						if key, err2 := strconv.Atoi(arr[2]); err2 == nil {
							d.is_exist = true
							d.RSAdecode.Numb = int64(n)
							d.RSAdecode.DecodeExp = int64(key)
						}
					}
				}
			} else if arr[0] == "text" {
				if key, err := strconv.Atoi(arr[1]); err == nil {				
					text := strings.Join(arr[1:], " ")
					text = EncodeText(text, int(d.RSAencode.Encode(int64(key))))
					fmt.Print(text)
				}
			}
		} else {
			continue		
		}

		fmt.Print(msg)
	}
}


func (d dialog) handleConnection() {
	go d.receiveMessage()
	time.Sleep(time.Second * 1)
	d.sendMessage()
}


func DecodeText(s string, key int) string {
	return s

	arr := []byte(s)
	for i := 0; i < len(arr); i++ {
		arr[i] -= byte(key)
	}
	return string(arr)
}


func EncodeText(s string, key int) string {
	return s
	
	arr := []byte(s)
	for i := 0; i < len(arr); i++ {
		arr[i] -= byte(key)
	}
	return string(arr)
}
