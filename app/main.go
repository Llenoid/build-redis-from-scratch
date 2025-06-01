package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on oprt :6379")

	// Create new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println("Error reading response: ", err)
		}

		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}
