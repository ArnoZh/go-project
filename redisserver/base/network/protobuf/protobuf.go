// Package protobuf .

package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"redisserver/base/network/msgpool"
	"redisserver/base/util"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// IMsg proto消息接口
type IMsg interface {
	proto.Message
	MarshalToSizedBuffer(dAtA []byte) (int, error)
	Size() int
}

// Processor 处理器
type Processor struct {
	littleEndian bool                    // 字节序
	id2type      map[uint32]reflect.Type // id to type映射
	type2id      map[reflect.Type]uint32 // type to id映射
	lenMsgID     int
	maxMsgLen    int
}

// NewProcessor 创建处理器
func NewProcessor(littleEndian bool, maxMsgLen, lenMsgID int) *Processor {
	processor := new(Processor)
	processor.littleEndian = littleEndian
	processor.type2id = make(map[reflect.Type]uint32)
	processor.id2type = make(map[uint32]reflect.Type)
	processor.maxMsgLen = maxMsgLen
	processor.lenMsgID = lenMsgID
	return processor
}

// SetByteOrder 设置字节序
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// Register 注册类型
func (p *Processor) Register(msg proto.Message) {
	msgType := reflect.TypeOf(msg)
	msgID := util.BKDRHash(msgType.Elem().Name())
	// 必须不为空 且为指针
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		logrus.Fatal("protobuf message pointer required")
	}
	// 不能重复注册
	if _, ok := p.type2id[msgType]; ok {
		logrus.Fatalf("message %s is already registered", msgType)
	}
	// 双向注册
	p.type2id[msgType] = msgID
	p.id2type[msgID] = msgType
}

// ReadMsgID 读取消息ID
func (p *Processor) ReadMsgID(data []byte) (uint32, error) {
	// 消息ID占4字节
	if len(data) < p.lenMsgID {
		return 0, errors.New("protobuf data too short")
	}
	// 根据字节序转化为对应的无符号整数
	var id uint32
	if p.littleEndian {
		id = binary.LittleEndian.Uint32(data)
	} else {
		id = binary.BigEndian.Uint32(data)
	}
	return id, nil
}

// Unmarshal 反序列化
func (p *Processor) Unmarshal(data []byte) (uint32, IMsg, error) {
	if len(data) < p.lenMsgID {
		return 0, nil, errors.New("protobuf data too short")
	}
	// 取消息ID
	var id uint32
	if p.littleEndian {
		id = binary.LittleEndian.Uint32(data)
	} else {
		id = binary.BigEndian.Uint32(data)
	}
	// 根据ID取类型
	typ, ok := p.id2type[id]
	if !ok {
		return 0, nil, fmt.Errorf("message id %v not registered", id)
	}
	// 根据类型反序列化
	msg0 := reflect.New(typ.Elem()).Interface().(IMsg)
	err := proto.Unmarshal(data[p.lenMsgID:], msg0)
	return id, msg0, err
}

// Marshal 序列化
func (p *Processor) Marshal(msg IMsg) ([]byte, error) {
	// 取实例类型
	msgType := reflect.TypeOf(msg)
	// 查找id
	_id, ok := p.type2id[msgType]
	if !ok {
		err := fmt.Errorf("msg %s not registered", msgType)
		return nil, err
	}
	// 序列化为字节流
	m, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	msglen := p.lenMsgID + len(m)
	if msglen == 0 || msglen > p.maxMsgLen {
		return nil, fmt.Errorf("msg %v too long: %d, max: %d", reflect.TypeOf(msg), msglen, p.maxMsgLen)
	}
	// 加入类型id作为消息头 完整消息
	data := make([]byte, p.lenMsgID, msglen)
	if p.littleEndian {
		binary.LittleEndian.PutUint32(data, _id)
	} else {
		binary.BigEndian.PutUint32(data, _id)
	}
	data = append(data, m...)
	return data, err
}

// MarshalToBuffer 序列化到指定buffer
func (p *Processor) MarshalToBuffer(msg IMsg, refCnt int32) (*msgpool.Buffer, error) {
	size := msg.Size()
	msglen := size + p.lenMsgID
	if msglen > p.maxMsgLen {
		return nil, fmt.Errorf("msg %v too long: %d, max: %d", reflect.TypeOf(msg), msglen, p.maxMsgLen)
	}
	// 取实例类型
	msgType := reflect.TypeOf(msg)
	// 查找id
	msgid, ok := p.type2id[msgType]
	if !ok {
		err := fmt.Errorf("msg %s not registered", msgType)
		return nil, err
	}
	buffer := msgpool.Get(msglen, refCnt)
	if buffer == nil {
		return nil, fmt.Errorf("can not get buffer from msgpool, type: %v size: %d", reflect.TypeOf(msg), msglen)
	}
	// 编码消息ID
	if p.littleEndian {
		binary.LittleEndian.PutUint32(buffer.Bytes[0:p.lenMsgID], msgid)
	} else {
		binary.BigEndian.PutUint32(buffer.Bytes[0:p.lenMsgID], msgid)
	}

	n, err := msg.MarshalToSizedBuffer(buffer.Bytes[p.lenMsgID:msglen])
	if err != nil {
		return buffer, err
	}
	buffer.Bytes = buffer.Bytes[:p.lenMsgID+n]
	return buffer, nil
}

// MsgID 根据类型取ID
func (p *Processor) MsgID(msg IMsg) uint32 {
	return p.type2id[reflect.TypeOf(msg)]
}

// MsgType 获取消息类型
func (p *Processor) MsgType(msgID uint32) reflect.Type {
	return p.id2type[msgID]
}

// Range 迭代调用
func (p *Processor) Range(f func(id uint32, t reflect.Type)) {
	for id, typ := range p.id2type {
		f(id, typ)
	}
}

// 获取消息名
func (p *Processor) MsgName(msgid uint32) string {
	for id, typ := range p.id2type {
		if msgid == id {
			return typ.Elem().Name()
		}
	}
	return ""
}
