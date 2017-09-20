package log

import (
	"sync"
	"io"
	"time"
	"runtime"
	"fmt"
)

// 日志等级
const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Llevel						  // log level
	LstdFlags     = Ldate | Ltime | Lmicroseconds | Llevel // initial values for the standard logger
)

// 日志类型
type Logger struct {
	mu     sync.Mutex 	// ensures atomic writes; protects the following fields
	prefix string     	// prefix to write at beginning of each line
	flag   int        	// properties
	out    []io.Writer  // destination for output
	buf    []byte     	// for accumulating text to write
	level	int		  	// 等级
	sep		string		// 字段分隔符
}

func (l *Logger) SetOutput(w []io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer  l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) SetFlag(flag int) {
	l.mu.Lock()
	defer  l.mu.Unlock()
	l.flag = flag
}

func (l *Logger) SetLevel(level int) {
	l.mu.Lock()
	defer  l.mu.Unlock()
	l.level = level
}

func (l *Logger) SetSeperator(sep string)  {
	l.mu.Lock()
	defer  l.mu.Unlock()
	l.sep = sep
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func loglevel(level int) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelInfo:
		return "INFO"
	case LogLevelFatal:
		return "FATAL"
	}
	return "UNKNOW"
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int, level int) {
	if len(l.prefix) > 0 {
		*buf = append(*buf, l.prefix...)
		*buf = append(*buf, l.sep...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '-')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '-')
			itoa(buf, day, 2)
			//*buf = append(*buf, l.sep...)
			if l.flag&(Ltime|Lmicroseconds) != 0 {
				*buf = append(*buf, ' ')
			}
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			//*buf = append(*buf, ' ')
			//*buf = append(*buf, l.sep...)
		}
		*buf = append(*buf, l.sep...)
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		//*buf = append(*buf, ": "...)
		*buf = append(*buf, l.sep...)
	}
	if l.flag&(Llevel) != 0 {
		*buf = append(*buf, loglevel(level)...)
		*buf = append(*buf, l.sep...)
	}
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(level int, calldepth int, s string) error {
	// Get time early if we need it.
	var now time.Time
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		now = time.Now()
	}
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, level)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	for _, out := range l.out {
		if out == nil {
			continue
		}
		_, err := out.Write(l.buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(level int, format string, v ...interface{}) {
	l.Output(level,2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(level int, v ...interface{}) { l.Output(level,2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Println(level int, v ...interface{}) { l.Output(level,2, fmt.Sprintln(v...)) }

func (l *Logger) Debug(v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.Println(LogLevelDebug, v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.Printf(LogLevelDebug, format, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.Println(LogLevelInfo, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.Printf(LogLevelInfo, format, v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.Println(LogLevelWarn, v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level <= LogLevelWarn {
		l.Printf(LogLevelWarn, format, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.level <= LogLevelError {
		l.Println(LogLevelError, v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.Printf(LogLevelError, format, v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level <= LogLevelFatal {
		l.Println(LogLevelFatal, v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.level <= LogLevelFatal {
		l.Printf(LogLevelFatal, format, v...)
	}
}

