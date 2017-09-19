package common

import (
	"io"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"gsf/msg"
	"errors"
)

type TCPMsg struct {
	msgId int			// 消息编号
	msg proto.Message	// protobuf msg
}

func CreateTCPMsgByMsgId(msgId int) proto.Message {
	switch msgId {
	case msg.Cmd_RegSvr_RegReq:
		return &msg.CmdRegSvrRegReq{}
	case msg.Cmd_RegSvr_RegNtf:
		return &msg.CmdRegSvrRegNtf{}
	default:
		return nil
	}
}

// 二进制消息解析
// header前4个字节，[msgid, msglen]
type MsgBinaryParser struct {
	header []byte 		// 字节头部
	body []byte			// 数据部分
}

func (mp *MsgBinaryParser) MakePacket() (interface{}, error) {
	var pktLen int = int(binary.BigEndian.Uint16(mp.header[0:2]))
	var msgId int = int(binary.BigEndian.Uint16(mp.header[2:4]))
	msg := CreateTCPMsgByMsgId(msgId)
	if msg == nil {
		return nil, errors.New("create tcp protobuf msg failed")
	}
	if pktLen > 0 {
		err := proto.Unmarshal(mp.body, msg)
		if err != nil {
			return nil, err
		}
	}
	return &TCPMsg{msgId: msgId, msg: msg}, nil
}

func (mp *MsgBinaryParser) Read(reader io.Reader) (interface{}, error) {
	// 读取头部
	n, err := reader.Read(mp.header)
	if n != len(mp.header) {
		return nil, err
	}

	var pktLen int = int(binary.BigEndian.Uint16(mp.header[0:2]))
	if pktLen == 0 {
		// 空消息
		return mp.MakePacket()
	}

	// 读取body
	mp.body = make([]byte, pktLen)
	n, err = reader.Read(mp.body)
	if n != pktLen {
		return nil, err
	}
	return mp.MakePacket()
}

func (mp *MsgBinaryParser) Write(msg interface{}, writer io.Writer) error {
	if tmsg, ok := msg.(*TCPMsg); ok {
		var head [4]byte
		binary.BigEndian.PutUint16(head[2:4], uint16(tmsg.msgId))
		data, err := proto.Marshal(tmsg.msg)
		if err != nil {
			return err
		}
		binary.BigEndian.PutUint16(head[0:2], uint16(len(data)))
		_, err = writer.Write(head[:])
		if err != nil {
			return err
		}
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}
