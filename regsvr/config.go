package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"gsf/log"
)

type Config struct {
	Daemon 		int 						`json:"daemon"`					// 是否后台启动
	ServerAddr 	string 						`json:"server_addr"`			// 服务器地址
	LogConf 	log.DateFileLoggerParam		`json:"log_conf"`				// 日志参数配置
}

var config = &Config{}
var logger *log.Logger = nil

func (c *Config) Init(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	logger, err = log.NewDateFileLogger(config.LogConf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) String() string {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("config json failed, err:", err)
	}
	return string(data)
}
