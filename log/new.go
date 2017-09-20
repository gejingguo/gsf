package log

import (
	"io"
	"os"
	"errors"
	"strconv"
)

var flagMap = map[string]int {
	"date": 			Ldate,
	"time": 			Ltime,
	"microseconds": 	Lmicroseconds,
	"longfile": 		Llongfile,
	"shortfile": 		Lshortfile,
	"utc": 				LUTC,
	"level": 			Llevel,
}

var levelMap = map[string]int {
	"debug": 	LogLevelDebug,
	"info": 	LogLevelInfo,
	"warn": 	LogLevelWarn,
	"error": 	LogLevelError,
	"fatal": 	LogLevelFatal,
}

var dateFormatMap = map[string]int {
	"day": DateFileFormatDayly,
	"hour": DateFileFormatHourly,
}

var Std = New([]io.Writer{os.Stdout}, "", LstdFlags, LogLevelInfo, "|")

// 输出配置
type OuterConfig struct {
	Type	string		`json:"type"`			// 类型，支持console, datefile, file
	Param 	[]string	`json:"param"`			// 参数，console(nil), datefile(file, day|hour, size), file(file)
}

// 配置参数
type LoggerConfig struct {
	Flag 	[]string	`json:"flag"`			// 标记, "
	Level	string		`json:"level"`			// 日志等级
	Sep 	string		`json:"sep"`			// 字段分隔符
	Outer 	[]OuterConfig	`json:"outer"`		// 目标输出列表
}

func New(out []io.Writer, prefix string, flag int, level int, sep string) *Logger {
	return &Logger{prefix:prefix, flag:flag, out:out, level:level, sep:sep}
}

func NewByConfig(config *LoggerConfig) (*Logger, error) {
	logFlag := 0
	for _, flagname := range config.Flag {
		flag := flagMap[flagname]
		logFlag |= flag
	}
	logLevel := levelMap[config.Level]
	logOuter := make([]io.Writer, 0)
	for _, outerConf := range config.Outer {
		outer, err := NewLogWriter(outerConf.Type, outerConf.Param)
		if err != nil {
			return nil, err
		}
		logOuter = append(logOuter, outer)
	}
	return New(logOuter, "", logFlag, logLevel, config.Sep), nil
}

func NewLogWriter(wtype string, param []string) (io.Writer, error) {
	switch wtype {
	case "console": 	return os.Stdout, nil
	case "file":
		if len(param) < 1 {
			return nil, errors.New("logger param invalid")
		}
		return NewFileWriter(param[0])
	case "datefile":
		if len(param) < 3 {
			return nil, errors.New("logger param invalid")
		}
		df := dateFormatMap[param[1]]
		if df == 0 {
			return nil, errors.New("logger param invalid")
		}
		size, err := strconv.Atoi(param[2])
		if err != nil {
			return nil, err
		}
		return NewDateFileWriter(param[0], df, size)
	default:
		return nil, errors.New("logger type invalid")
	}
}
