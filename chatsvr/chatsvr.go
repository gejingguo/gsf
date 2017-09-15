package main

import (
	"runtime"
	"io"
	"bufio"
	"fmt"
	"os"

	"gsf/tcp"
)

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

type MsgParse struct {
}

func (mp MsgParse) Read(reader io.Reader) (interface{}, error) {
	br := bufio.NewReader(reader)
	line, _, err := br.ReadLine()
	if err == io.EOF && line != nil{
		return line, nil
	}
	return line, err
}

func (mp MsgParse) Write(msg interface{}, writer io.Writer) error {
	if s,ok := msg.([]byte); ok {
		s = append(s, '\n')
		_,err := writer.Write(s)
		return err
	}
	return nil
}

// 用户
type User struct {
	session *tcp.TCPSession		// 对应会话
	name string 				// 名称
}
var UserMap = make(map[string]User)

func main() {

	listener, err := tcp.NewTCPListener(":9001", func(session *tcp.TCPSession) {
		session.MsgParse = MsgParse{}
		session.ErrorHandler = func(session *tcp.TCPSession, err error) {
			fmt.Println("session closed,", session.Id, err)
		}
		session.MsgHandler = func(session *tcp.TCPSession, msg interface{}) {
			if s,ok := msg.([]byte); ok {
				fmt.Println("session recv:", session.Id, s)
				// session.Write(msg)
				for _, s := range tcp.SessionMap {
					s.Write(msg)
				}
			}
		}
		session.Start()
	})
	listener.Start()

	if err != nil {
		fmt.Println(err)
	}
	tcp.Run()
}

