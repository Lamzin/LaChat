package main


import (
	// "fmt"
	// "net"
	// "bufio"
	// "os"
	// "./rsa"
	// "./connection"
	"./msghnd"
)


func main() {
	var conn msghnd.ConnClient
	conn.Connect("127.0.0.1:6000")
	conn.HandlerIO()
}
