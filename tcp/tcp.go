package tcp

import (
	"net"
	"fmt"
	"runtime"
	"os"
	"time"
	"io"
	"errors"
)

func printStack() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

type TCPAcceptHandler func(*TCPSession)
// Session错误处理
type TCPSessionErrorHandler func(session *TCPSession, err error)
// Session消息处理
type TCPSessionMsgHandler func(session *TCPSession, msg interface{})
// session connect连接成功处理
type TCPSessionConnectedHandler func (session *TCPSession, err error)
//
type TCPTickHandler func ()
//
type TCPListener struct {
	listener net.Listener
	acceptf TCPAcceptHandler
}

type TCPSession struct {
	Id uint64 								// 唯一编号
	Conn net.Conn							// 连接
	ConnectAddr string 						// 连接服务器地址
	MsgParse ProtocolParse 					// 消息解析接口
	ErrorHandler TCPSessionErrorHandler 	// 错误处理
	MsgHandler TCPSessionMsgHandler			// 消息处理
	ConnectHandler TCPSessionConnectedHandler // connect处理
	SendMsgChan chan interface{}			// 消息发送队列
	UserData interface{}					// 用户数据
}

type TCPAcceptData struct {
	conn net.Conn
	acceptf TCPAcceptHandler
}

type TCPConnectData struct {
	session *TCPSession
	conn net.Conn
	err error
}

type TCPMsg struct {
	Session *TCPSession
	Object interface{}
	Err error
}

// ProtocolParse 协议消息解析封包接口
// 该接口操作执行在读写go程中，不在主逻辑线程注意
type ProtocolParse interface {
	Read(reader io.Reader) (interface{}, error)
	Write(object interface{}, writer io.Writer) error
}

// 环境
type TCPContext struct {
	sessionId uint64
	connectChan chan TCPConnectData
	acceptChan chan TCPAcceptData
	msgChan chan TCPMsg
	tickChan <-chan time.Time
	tickHandler TCPTickHandler
}
var context = TCPContext{connectChan: make(chan TCPConnectData, 256), acceptChan: make(chan TCPAcceptData, 1024), msgChan: make(chan TCPMsg, 1024)}
// 所有已经连接的会话
var SessionMap = make(map[uint64]*TCPSession)

// 开始监听
func (l *TCPListener) Start() {
	go func(listener net.Listener, acceptf TCPAcceptHandler, acceptChan chan TCPAcceptData) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("ln accept fail.", err)
				continue
			}
			fmt.Println("goaccept addr;", conn.RemoteAddr())
			acceptChan <- TCPAcceptData{conn, acceptf}
		}
	}(l.listener, l.acceptf, context.acceptChan)
}

// NewTCPListener 创建新的TCP监听器
func NewTCPListener(addr string, handler TCPAcceptHandler) (*TCPListener, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TCPListener{listener: ln, acceptf:handler}, nil
}

// 主动关闭
var ErrClosedActive = errors.New("closed active")

// 主动关闭
func (s *TCPSession) CloseByActive() {
	s.Close(ErrClosedActive)
}

func (s *TCPSession) IsConnected() bool {
	return s.Conn != nil
}

func (s *TCPSession) IsConnectedSession() bool {
	return s.ConnectAddr != ""
}

//
func (s *TCPSession) Close(err error) {
	if s.Conn != nil {
		s.Conn.Close()
	}
	if s.ErrorHandler != nil {
		s.ErrorHandler(s, err)
	}
	if s.SendMsgChan != nil {
		close(s.SendMsgChan)
		s.SendMsgChan = nil
	}
	// 从会话列表中删除
	delete(SessionMap, s.Id)
}

// 创建连接其他服务器的会话
func NewConnectTCPSession(addr string, handler TCPSessionConnectedHandler) *TCPSession {
	context.sessionId++
	session := &TCPSession{Id: context.sessionId, ConnectAddr:addr, ConnectHandler: handler}
	session.SendMsgChan = make(chan interface{}, 1024)
	// 加入列表
	SessionMap[session.Id] = session

	// 拨号建立连接
	go func(s *TCPSession, connChan chan TCPConnectData) {
		conn, err := net.Dial("tcp", s.ConnectAddr)
		connChan <- TCPConnectData{s, conn, err}
	}(session, context.connectChan)

	return session
}

// 开始读写
func (s *TCPSession) Start() {
	if s.Conn == nil {
		return
	}
	go func(s *TCPSession) {
		for {
			if s.MsgParse != nil {
				obj, err := s.MsgParse.Read(s.Conn)
				context.msgChan<- TCPMsg{s,obj, err}
				if err != nil {
					break
				}
			}
		}
		fmt.Println("conn read go finised.")
	}(s)
	go func(s *TCPSession) {
		for {
			if s.MsgParse != nil {
				msg, ok := <- s.SendMsgChan
				if !ok {
					break
				}
				err := s.MsgParse.Write(msg, s.Conn)
				if err != nil {
					// 写入错误，通知逻辑层关闭会话
					context.msgChan <- TCPMsg{s, nil, err}
					break
				}
			}
		}
		fmt.Println("conn write go finished.")
	}(s)
}

// 写入数据异步
func (s *TCPSession) Write(msg interface{}) bool{
	select {
	case s.SendMsgChan<- msg:
		return true
	default:
		return false
	}
}

func SetTick(ms time.Duration, update TCPTickHandler) {
	context.tickChan = time.Tick(ms*time.Millisecond)
	context.tickHandler = update
}

func Run() {
	for {
		select {
		case acceptData := <-context.acceptChan:
			context.sessionId++
			session := &TCPSession{Id: context.sessionId, Conn: acceptData.conn}
			session.SendMsgChan = make(chan interface{}, 1024)
			// 加入列表
			SessionMap[session.Id] = session
			acceptData.acceptf(session)
			// 消息队列
		case msg := <-context.msgChan:
			if msg.Err != nil {
				msg.Session.Close(msg.Err)
			} else {
				if msg.Session.MsgHandler != nil {
					msg.Session.MsgHandler(msg.Session, msg.Object)
				}
			}
		case connData := <- context.connectChan:
			if connData.session != nil {
				if connData.err == nil {
					connData.session.Conn = connData.conn
				}
				if connData.session.ConnectHandler != nil {
					connData.session.ConnectHandler(connData.session, connData.err)
				}
			}
		case <- context.tickChan:
			if context.tickHandler != nil {
				context.tickHandler()
			}
		}
	}
}


