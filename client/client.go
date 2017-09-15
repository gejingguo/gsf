package main

import (
	"gsf/tcp"
	"fmt"
)

func main() {

	tcp.NewConnectTCPSession("180.97.33.108:80", func(session *tcp.TCPSession, err error) {
		if err == nil {
			fmt.Println("connect ok.")
		} else {
			fmt.Println("connect failed.")
			session.Close(err)
		}
	})

	tcp.Run(1000, func() {

	})
}
