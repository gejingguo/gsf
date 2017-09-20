package main

import (
	"net"
	"fmt"
	"gsf/msg"
	"gsf/common"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9002")
	if err != nil {
		fmt.Println("connect server failed, err:", err)
		return
	}

	req := &msg.CmdRegSvrRegReq {
		Server: &msg.ServerInfo{
			Id: 1,
			Type: 1,
			Group: 1,
			ListenAddr: "localhost:9002",
		},
	}

	msgParser := common.MsgBinaryParser{}
	tmsg := &common.TCPMsg{
		MsgId: int(msg.Cmd_RegSvr_RegReq),
		Msg: req,
	}
	err = msgParser.Write(tmsg, conn)
	if err != nil {
		fmt.Println("write msg failed, err:", err)
		return
	}
	fmt.Println("write msg ok.")

	for {
		time.Sleep(1)
	}
}
