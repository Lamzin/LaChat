package main


import (
	"fmt"
	"net"
	"os"
	"./msghnd"
)


func main() {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}


	// var conn_arr []msghnd.ConnServer

	// for {
		var conn msghnd.ConnServer
		if connection, err := ln.Accept(); err != nil {
			fmt.Print("error\n\n")
			fmt.Print(err)
			// continue
			os.Exit(1)
		} else {
			conn.Conn = connection
			// conn_arr = append(conn_arr, conn)
			fmt.Print(connection.RemoteAddr())
		}

		// go conn.HandlerIO()
		conn.HandlerIO()
	// }

}
