package main

import (
	"github.com/golang/protobuf/proto"
	"fmt"
	"gsf/log"
)

func main() {
	/*
	tcp.NewConnectTCPSession("180.97.33.108:80", func(session *tcp.TCPSession, err error) {
		if err == nil {
			fmt.Println("connect ok.")
		} else {
			fmt.Println("connect failed.")
			session.Close(err)
		}
	})

	tcp.Run()
	*/
	proto.Marshal()

	/*logw, err := log.NewDateFileWriter("D:\\client.log", log.DateFileFormatDayly, 10000)
	if err != nil {
		fmt.Println("open faile failed, err:", err)
		return
	}
	logger := log.New(logw, "", log.LstdFlags, log.LogLevelDebug, "|")*/
	logger, err := log.NewDateFileLogger(log.DateFileLoggerParam{
		Flag: log.LstdFlags,
		Level: log.LogLevelDebug,
		Sep: "|",
		File: "D:\\client.log",
		Df: log.DateFileFormatDayly,
		Size: 10240,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	logger.Debugf("pb: %v", tab)
	logger.Infof("pb: %v", tab)
	logger.Warnf("pb: %v", tab)
	logger.Errorf("pb: %v", tab)
	logger.Fatalf("pb: %v", tab)
}
