// Package chanrpc .

package chanrpc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"redisserver/base/util/cbctx"
)

// SnapshotMarker marker 消息
type SnapshotMarker struct {
}

// DBLoadClean 清除标签，play与DB消息通道中的异步加载消息
type DBLoadClean struct {
}

// CallInfoObj 需要落地的调用参数
type CallInfoObj struct {
	SerialNo int64  // 保存消息序号
	Owner    string // 消息所属模块名

	ID      int64    // 消息类型id
	ReqType string   // 消息类型
	Req     bson.Raw // 入参
	Ctx     cbctx.M  // 回调上下文
	HasRet  bool     // 是否已经返回 由被调用方使用
	Time    string
}

// RetInfoObj 需要落地的结果消息
type RetInfoObj struct {
	SerialNo int64  // 保存消息序号
	ID       int64  // 消息Ack类型id
	Owner    string // 消息所属模块名

	AckType string   // 消息类型
	Ack     bson.Raw // 结果值 作为回调函数的入参
	ErrStr  string   // 错误
	Ctx     cbctx.M  // 回调上下文
	Time    string
}

type msgType struct {
	typ  reflect.Type
	kind reflect.Kind
}

// MessageIDMap 消息ID对应Type
var messageIDMap = map[uint32]msgType{}

// registerMsgIDType 注册消息ID对应类型
func registerMsgIDType(msgID uint32, msg interface{}) {

	typ := reflect.TypeOf(msg)

	kind := typ.Kind()
	if kind == reflect.Ptr {
		typ = typ.Elem()
	}

	messageIDMap[msgID] = msgType{
		typ:  typ,
		kind: kind,
	}
}

// UnmarshalBSON 将bson解码为结果信息
func (ri *RetInfo) UnmarshalBSON(data []byte) error {
	riObj := &RetInfoObj{}
	err := bson.Unmarshal(data, riObj)
	if err != nil {
		return err
	}

	ri.Ctx = riObj.Ctx
	ri.SerialNo = uint32(riObj.SerialNo)

	msgType, ok := messageIDMap[uint32(riObj.ID)]
	if !ok {
		return fmt.Errorf("req %d typ not found", riObj.ID)
	}

	iAck := reflect.New(msgType.typ).Interface()
	err = bson.Unmarshal(riObj.Ack, iAck)
	if err != nil {
		return err
	}

	if msgType.kind == reflect.Struct {
		iAck = reflect.ValueOf(iAck).Elem()
	}

	ri.Owner = riObj.Owner
	ri.Ack = iAck

	if riObj.ErrStr != "" {
		ri.Err = fmt.Errorf(riObj.ErrStr)
	}

	return nil
}

// MarshalBSON 结果信息编码为bson结构
func (ri *RetInfo) MarshalBSON() ([]byte, error) {
	riObj := &RetInfoObj{
		SerialNo: int64(ri.SerialNo),
		Owner:    ri.Owner,
		ID:       int64(ri.GetMsgID()),
		Ctx:      ri.Ctx,
		Time:     time.Now().UTC().String(),
	}
	b0, err := bson.Marshal(ri.Ack)
	if err != nil {
		logrus.Errorf("marshal ack: %s, ack %#v", err.Error(), ri.Ack)
		return nil, err
	}
	if ri.Ack != nil {
		riObj.AckType = reflect.TypeOf(ri.Ack).Elem().Name()
	}
	riObj.Ack = b0

	if ri.Err != nil {
		riObj.ErrStr = ri.Err.Error()
	}

	b2, err := bson.Marshal(riObj)
	if err != nil {
		logrus.Errorf("xxxx riObj error: %s", err.Error())
		return nil, err
	}
	return b2, nil
}

// UnmarshalBSON 将bson解码为调用参数
func (ci *CallInfo) UnmarshalBSON(data []byte) error {
	ciObj := &CallInfoObj{}
	err := bson.Unmarshal(data, ciObj)
	if err != nil {
		return err
	}

	ci.id = uint32(ciObj.ID)
	ci.ctx = ciObj.Ctx
	ci.hasRet = ciObj.HasRet
	ci.Owner = ciObj.Owner
	ci.SerialNo = uint32(ciObj.SerialNo)

	msgType, ok := messageIDMap[uint32(ciObj.ID)]
	if !ok {
		return fmt.Errorf("req %d typ not found", ciObj.ID)
	}

	iReq := reflect.New(msgType.typ).Interface()
	err = bson.Unmarshal(ciObj.Req, iReq)
	if err != nil {
		return err
	}

	if msgType.kind == reflect.Struct {
		iReq = reflect.ValueOf(iReq).Elem()
	}
	ci.Req = iReq
	return nil
}

// MarshalBSON 调用参数编码为bson结构
func (ci *CallInfo) MarshalBSON() ([]byte, error) {
	ciObj := &CallInfoObj{
		SerialNo: int64(ci.SerialNo),
		Owner:    ci.Owner,
		ID:       int64(ci.id),
		Ctx:      ci.ctx,
		HasRet:   ci.hasRet,
		Time:     time.Now().UTC().String(),
	}
	b, err := bson.Marshal(ci.Req)
	if err != nil {
		return nil, err
	}

	if ci.Req != nil {
		ciObj.ReqType = reflect.TypeOf(ci.Req).Elem().Name()
	}

	ciObj.Req = b

	b2, err := bson.Marshal(ciObj)
	if err != nil {
		return nil, err
	}
	return b2, nil
}

// IsMarker 是否是Marker消息
func (ci *CallInfo) IsMarker() bool {

	_, ok := ci.Req.(*SnapshotMarker)

	return ok
}

// IsAsyncClean 是否是Marker消息
func (ri *RetInfo) IsAsyncClean() bool {

	_, ok := ri.Ack.(*DBLoadClean)

	return ok
}

// CallInfoSort 调用信息排序
type CallInfoSort []*CallInfo

// Len 长度
func (ci CallInfoSort) Len() int {
	return len(ci)
}

// Swap 交换i, j
func (ci CallInfoSort) Swap(i, j int) {
	ci[i], ci[j] = ci[j], ci[i]
}

// Less elem(i) < elem(j)
func (ci CallInfoSort) Less(i, j int) bool {
	return ci[i].SerialNo < ci[j].SerialNo
}

// RetInfoSort 异步调用返回排序
type RetInfoSort []*RetInfo

// Len 长度
func (ri RetInfoSort) Len() int {
	return len(ri)
}

// Swap 交换i, j
func (ri RetInfoSort) Swap(i, j int) {
	ri[i], ri[j] = ri[j], ri[i]
}

// Less elem(i) < elem(j)
func (ri RetInfoSort) Less(i, j int) bool {
	return ri[i].SerialNo < ri[j].SerialNo
}
