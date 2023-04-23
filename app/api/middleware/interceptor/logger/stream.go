package logger

import (
	"github.com/hphphp123321/mahjong-server/app/global"
	"google.golang.org/grpc"
)

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		wrapper := newStreamServer(ss)
		return handler(srv, wrapper)
	}
}

// 嵌入式EdgeServerStream允许我们访问RecvMsg函数
type streamServer struct {
	grpc.ServerStream
}

func newStreamServer(s grpc.ServerStream) grpc.ServerStream {
	return &streamServer{s}
}

// RecvMsg从流中接收消息
func (e *streamServer) RecvMsg(m interface{}) error {
	// 在这里，我们可以对接收到的消息执行额外的逻辑，例如
	// 验证
	global.Log.Infof("Receive a message (Type: %T)\n", m)
	if err := e.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func (e *streamServer) SendMsg(m interface{}) error {
	global.Log.Debugf("Send a message (Type: %T)\n", m)
	if err := e.ServerStream.SendMsg(m); err != nil {
		return err
	}
	return nil
}
