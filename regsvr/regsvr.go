package main

import (
	"os"
	"fmt"
	"gsf/tcp"
	"gsf/common"
)

func usage() {
	fmt.Println(os.Args[0], " configfile")
	os.Exit(-1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	err := config.Init(os.Args[1])
	if err != nil {
		fmt.Println("config init failed, err:", err)
		os.Exit(-1)
	}

	listener, err := tcp.NewTCPListener( config.ServerAddr, func(session *tcp.TCPSession) {
		session.MsgParse = &common.MsgBinaryParser{}
		session.ErrorHandler = func(session *tcp.TCPSession, err error) {
			//fmt.Println("session closed,", session.Id, err)
			logger.Infof("session(%v) closed, err:%v", session.Id, err)
		}
		session.MsgHandler = func(session *tcp.TCPSession, msg interface{}) {
			fmt.Println("conn recv msg:", msg)
			if s,ok := msg.(*common.TCPMsg); ok {
				logger.Info("recv msg:", s.Msg)
			}
		}
		session.Start()
		fmt.Println("recv new conn.")
	})
	if err != nil {
		logger.Error("tcplistener create failed, err:", err)
		return
	}
	listener.Start()
	logger.Info("server start ok.")
	fmt.Println("server start ok...")

	tcp.Run()
}
