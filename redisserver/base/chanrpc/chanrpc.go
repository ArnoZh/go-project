// Package chanrpc .

package chanrpc

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"redisserver/base/conf"
	"redisserver/base/util"
	"redisserver/base/util/cbctx"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
)

// Router 消息路由器
type Router interface {
	Route(id uint32, x proto.Message) error
}

// Handler 方法句柄  处理CallInfo
type Handler func(ci *CallInfo)

// Callback 回调
type Callback func(ri *RetInfo)

// Server 代理服务器
type Server struct {
	functions map[interface{}]Handler // 处理func
	ChanCall  chan *CallInfo          // call时往chan传入值
	// router    Router
}

// CallInfo 调用参数
type CallInfo struct {
	id      uint32        // 消息类型id
	Req     interface{}   // 入参
	chanRet chan *RetInfo // 结果信息返回通道
	cb      Callback      // 回调 todo 取消
	ctx     cbctx.M       // 回调上下文
	hasRet  bool          // 是否已经返回 由被调用方使用

	Owner    string // 消息所属模块
	SerialNo uint32 // 停服时保存消息序号
}

// Ret 调用 请求的 回调
func (ci *CallInfo) Ret(ret interface{}) {
	// 检查回调是否已经使用过
	if ci.hasRet {
		logrus.Errorf("chanrpc can not ret twice, %v", string(debug.Stack()))
		return
	}
	// 标记
	ci.hasRet = true
	// 封装参数 执行回调
	err := ci.ret(&RetInfo{Ack: ret, Ctx: ci.ctx})
	if err != nil {
		logrus.Errorf("chanrpc ret error: %v, msgid: %v", err, ci.id)
	}
}

// Ack 调用请求的回调
func (ci *CallInfo) ret(ri *RetInfo) (err error) {
	// 检查返回通道
	if ci.chanRet == nil {
		return
	}
	// 错误捕获
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	// 封装参数 将结果信息放入返回通道
	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

// GetMsgID 调用消息ID
func (ci *CallInfo) GetMsgID() uint32 {
	return ci.id
}

// GetMsgName 返回消息的结果(回调入参)类型名
func (ci *CallInfo) GetMsgName() string {
	if ci.Req == nil {
		return "UnknownCall"
	}
	return MsgName(ci.Req)
}

// RetInfo 结果信息
type RetInfo struct {
	Ack interface{} // 结果值 作为回调函数的入参
	Err error       // 错误
	cb  Callback    // 回调 todo 取消
	Ctx cbctx.M     // 回调上下文

	SerialNo uint32 // 停服时保存消息序号
	Owner    string // 接收消息的模块
}

// GetMsgID 返回消息的结果(回调入参)类型ID
func (ri *RetInfo) GetMsgID() uint32 {
	if ri.Err != nil || ri.Ack == nil {
		return 0
	}
	return MsgID(ri.Ack)
}

// GetMsgName 返回消息的结果(回调入参)类型名
func (ri *RetInfo) GetMsgName() string {
	if ri.Ack == nil {
		unknown := "UnknownRet"
		if ri.Err != nil {
			unknown += ri.Err.Error()
		}
		return unknown
	}
	return MsgName(ri.Ack)
}

// Client 客户端
type Client struct {
	chanCall        chan *CallInfo           // 调用信息通道
	chanSyncRet     chan *RetInfo            // 同步调用结果通道
	ChanAsynRet     chan *RetInfo            // 异步调用结果通道
	pendingAsynCall int                      // 最大待处理异步调用
	functions       map[interface{}]Callback // ack 对应处理函数
	mylock          sync.Mutex               // 加锁
}

// AsyncRet 异步调用结果
func (c *Client) AsyncRet() chan *RetInfo {
	return c.ChanAsynRet
}

// IMsgID 消息可实现该接口来自定义MsgID，达成如消息结构体复用等高级功能
type IMsgID interface {
	MsgID() uint32
}

// MsgID 求消息的消息ID，传入值必须是指针
func MsgID(m interface{}) uint32 {
	if msgIDGen, ok := m.(IMsgID); ok {
		return msgIDGen.MsgID()
	}
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Struct {
		return util.BKDRHash(typ.Name())
	}
	return util.BKDRHash(typ.Elem().Name())
}

// MsgName 求消息的消息名
func MsgName(m interface{}) string {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Struct {
		return typ.Name()
	}
	return typ.Elem().Name()
}

// NewServer 新建服务器
func NewServer(l int) *Server {
	s := new(Server)
	s.functions = make(map[interface{}]Handler)
	s.ChanCall = make(chan *CallInfo, l)
	return s
}

// Register 向服务器注册处理函数 根据id索引
func (s *Server) Register(msg interface{}, f Handler) {
	msgID := MsgID(msg)
	//logrus.Infof("%v:%v", msgID, MsgName(msg))
	if _, ok := s.functions[msgID]; ok {
		panic(fmt.Sprintf("function ID %v: already registered", reflect.TypeOf(msg)))
	}
	s.functions[msgID] = f
	registerMsgIDType(msgID, msg)
}

// exec 实际执行
func (s *Server) exec(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				err = fmt.Errorf("%v: %s", r, buf[:l])
			} else {
				err = fmt.Errorf("%v", r)
			}
			if !ci.hasRet {
				// 如果是异步投递，那么消息不会被返回
				_ = ci.ret(&RetInfo{
					Err: fmt.Errorf("%v", r),
					Ctx: ci.ctx,
				})
			}
		}
	}()
	// 根据id取handlenot register handlerr
	handler, ok := s.functions[ci.id]
	if !ok {
		// panic will be defer and call ci.Ret()
		panic(fmt.Sprintf("msg %+v %v not register handler", ci.id, ci.Req))
	}
	handler(ci)
	return
}

// Exec 执行
func (s *Server) Exec(ci *CallInfo) {
	ci.hasRet = false
	err := s.exec(ci)
	if err != nil {
		logrus.Errorf("%v", err)
	}
}

// Cast 直接投递消息  忽略任何错误和返回值
func (s *Server) Cast(req interface{}) (err error) {
	id := MsgID(req)
	err = call(
		s.ChanCall,
		&CallInfo{
			id:  id,
			Req: req,
		},
		false,
	)
	return
}

// CastMustSuccess 直接投递消息 阻塞方式，必须成功
func (s *Server) CastMustSuccess(req interface{}) (err error) {
	id := MsgID(req)
	err = call(
		s.ChanCall,
		&CallInfo{
			id:  id,
			Req: req,
		},
		true,
	)
	return
}

// Call 启动一个client来进行调用
func (s *Server) Call(req interface{}) *RetInfo {
	return s.Open(0).Call(req)
}

// Close 关闭服务器
func (s *Server) Close() {
	close(s.ChanCall)
	for ci := range s.ChanCall {
		_ = ci.ret(&RetInfo{
			Err: errors.New("chanrpc server closed"),
		})
	}
}

// Open 启动一个客户端
func (s *Server) Open(l int) *Client {
	c := NewClient(l)
	c.Attach(s)
	return c
}

// NewClient 新建客户端 设置异步结果通道的最大缓冲值
func NewClient(l int) *Client {
	c := &Client{
		chanSyncRet: make(chan *RetInfo, 1),
		ChanAsynRet: make(chan *RetInfo, l),
		functions:   make(map[interface{}]Callback),
		mylock:      sync.Mutex{},
	}
	return c
}

// AckRegister 向服务器注册Ack处理函数
func (c *Client) AckRegister(msg interface{}, f Callback) {
	msgID := MsgID(msg)
	if _, ok := c.functions[msgID]; ok {
		panic(fmt.Sprintf("function ID %v: already registered", reflect.TypeOf(msg)))
	}
	c.functions[msgID] = f
	registerMsgIDType(msgID, msg)
}

// Attach 将client的请求通道依附于服务器
func (c *Client) Attach(s *Server) {
	c.chanCall = s.ChanCall
}

// Call0 同步无结果调用
func (c *Client) Call0(id uint32, req interface{}) error {
	err := call(c.chanCall, &CallInfo{
		id:      id,
		Req:     req,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return err
	}
	ri := <-c.chanSyncRet
	return ri.Err
}

// Call 同步有结果调用
func (c *Client) Call(req interface{}) *RetInfo {
	id := MsgID(req)
	err := call(c.chanCall, &CallInfo{
		id:      id,
		Req:     req,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return &RetInfo{
			Err: err,
		}
	}
	ri := <-c.chanSyncRet
	return ri
}

// AsynCall 异步调用
func (c *Client) AsynCall(req interface{}, cb Callback, ctx cbctx.M) error {
	defer c.mylock.Unlock()
	c.mylock.Lock()
	id := MsgID(req)
	err := call(c.chanCall, &CallInfo{
		id:      id,
		Req:     req,
		chanRet: c.ChanAsynRet,
		cb:      cb,
		ctx:     ctx,
	}, false)
	if err != nil {
		return err
	}
	c.pendingAsynCall++
	return nil
}

// execCb 执行回调
func execCb(ri *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				logrus.Errorf("%v: %s", r, debug.Stack())
			} else {
				logrus.Errorf("%v", r)
			}
		}
	}()
	ri.cb(ri)
}

// Cb 执行回调
func (c *Client) Cb(ri *RetInfo) {

	c.mylock.Lock()
	c.pendingAsynCall--

	// 落地的消息，cb为空
	if ri.cb == nil {
		msgID := ri.GetMsgID()
		ri.cb = c.functions[msgID]
	}
	c.mylock.Unlock() //在这就解锁，防止回调里面死锁
	execCb(ri)
}

// Close 关闭client
func (c *Client) Close() {

	c.mylock.Lock()
	nNumber := c.pendingAsynCall
	c.mylock.Unlock() //防止回调函数中造成死锁
	for nNumber > 0 {
		c.mylock.Lock()
		nNumber = c.pendingAsynCall
		c.mylock.Unlock()
		c.Cb(<-c.ChanAsynRet)
	}
}

// Idle 判断是否空闲
func (c *Client) Idle() bool {
	defer c.mylock.Unlock()
	c.mylock.Lock()
	return c.pendingAsynCall == 0
}

// call 调用,分阻塞与非阻塞模式,仅仅将请求放入请求通道
func call(chanCall chan *CallInfo, ci *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	if block { // 阻塞模式
		chanCall <- ci
		return
	}
	// 非阻塞模式
	select {
	case chanCall <- ci:
	default:
		return fmt.Errorf("server chanrpc channel full, msg %v %+v", reflect.TypeOf(ci.Req).String(), ci.Req)
	}
	return
}
