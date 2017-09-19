package main

import (
	"gsf/proto"
	gproto "github.com/golang/protobuf/proto"
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

	ab := &proto.AddressBook{
		People: []*proto.Person{
			{Name: "name1", Id:1, Email: "name1@qq.com"},
			{Name: "name2", Id:2, Email: "name2@qq.com"},
			{Name: "name3", Id:3, Email: "name3@qq.com"},
		},
	}

	data, err := gproto.Marshal(ab)
	if err != nil {
		fmt.Println("ab marshal failed, err:", err)
		return
	}

	tab := &proto.AddressBook{}
	err = gproto.Unmarshal(data, tab)
	if err != nil {
		fmt.Println("tab unmarshal failed, err:", err)
		return
	}
	//fmt.Println("pb info")
	//fmt.Println(tab)
	log.Std.Debug("pb info")
	log.Std.Debugf("pb: %v", tab)
	log.Std.Infof("pb: %v", tab)
	log.Std.Warnf("pb: %v", tab)
	log.Std.Errorf("pb: %v", tab)
	log.Std.Fatalf("pb: %v", tab)

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
	logger.Debugf("pb: %v", tab)
	logger.Infof("pb: %v", tab)
	logger.Warnf("pb: %v", tab)
	logger.Errorf("pb: %v", tab)
	logger.Fatalf("pb: %v", tab)
}
