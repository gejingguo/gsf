package log

import (
	"os"
	"time"
)

// 简单文件输出
// 同一个文件一直输出，不切文件
type FileWriter struct {
	filename 	string		// 文件名
	file 		*os.File
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func NewFileWriter(filename string) (*FileWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 667)
	if err != nil {
		return nil, err
	}
	return &FileWriter{filename: filename, file: file}, nil
}

const (
	DateFileFormatNone = iota
	DateFileFormatDayly
	DateFileFormatHourly
)

// 日期文件输出
// 支持按照日期进行切换，支持按照自然天，小时
// 支持按照文件大小切换
type DateFileWriter struct {
	fileName 		string 			// 文件名
	dateFormat 		int 			// 日期格式
	fileSize 		int				// 文件大小

	file 			*os.File
	lastTime 		time.Time		// 上次创建文件时间
	writeSize 		int				// 已经输出字节数量
	lastFileSuffix	string			// 上次切换文件后缀
}

func NewDateFileWriter(filename string, df int, size int) (*DateFileWriter, error) {
	suffix := newFileSuffix(df)
	if suffix != "" {
		filename += "." + suffix
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 667)
	if err != nil {
		return nil, err
	}

	return &DateFileWriter{fileName: filename, dateFormat: df, fileSize: size, file: file, lastTime: time.Now()}, nil
}

type DateFileLoggerParam struct{
	Flag 	int			// 标记
	Level	int			// 日志等级
	Sep 	string		// 字段分隔符
	File 	string		// 日志文件名含路径
	Df 		int 		// 日期格式，每日，还是每小时
	Size 	int 		// 日志单个文件大小上限
}
func NewDateFileLogger(p DateFileLoggerParam) (*Logger, error) {
	out, err := NewDateFileWriter(p.File, p.Df, p.Size)
	if err != nil {
		return nil, err
	}
	return New(out, "", p.Flag, p.Level, p.Sep), nil
}

func (w *DateFileWriter) Write(p []byte) (n int, err error) {
	n, err = w.file.Write(p)
	w.writeSize += n
	if w.checkNeedSwitchFile() {
		w.switchFile()
	}
	return
}

// 检测是否需要切换文件
func (w *DateFileWriter) checkNeedSwitchFile() bool {
	if w.fileSize > 0 && w.writeSize >= w.fileSize {
		return true
	}
	now := time.Now()
	if w.dateFormat == DateFileFormatDayly {
		if now.Year() != w.lastTime.Year() || now.Month() != w.lastTime.Month() || now.Day() != w.lastTime.Day() {
			return true
		}
	} else if w.dateFormat == DateFileFormatHourly {
		if now.Year() != w.lastTime.Year() || now.Month() != w.lastTime.Month() || now.Day() != w.lastTime.Day() || now.Hour() != w.lastTime.Hour() {
			return true
		}
	}
	return false
}

func newFileSuffix(df int) string {
	if df == DateFileFormatDayly {
		return time.Now().Format("2006-01-02.150405")
	} else if df == DateFileFormatHourly {
		return time.Now().Format("2006-01-02-15.0405")
	}
	return ""
}

func (w *DateFileWriter) newFileSuffix() string {
	suffix := newFileSuffix(w.dateFormat)
	if suffix == "" {
		suffix = w.lastFileSuffix
	}
	return suffix
}

// 切换文件
// 新建文件名称为xx.log.YY-MM-DD-HH
func (w *DateFileWriter) switchFile() {
	newFileSuffix := w.newFileSuffix()
	if newFileSuffix == w.lastFileSuffix {
		return
	}

	w.file.Close()
	w.lastFileSuffix = newFileSuffix
	w.lastTime = time.Now()
	w.writeSize = 0

	nf, err := os.OpenFile(w.fileName + "." + w.lastFileSuffix, os.O_APPEND|os.O_CREATE, 667)
	if err != nil {
		panic(err)
		return
	}
	w.file = nf
}
