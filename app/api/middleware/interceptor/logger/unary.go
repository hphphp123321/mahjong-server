package logger

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-server/app/global"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"strings"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		_, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			global.Log.Debugf("context has no metadata")
		}
		ip, err := getClientIP(ctx)
		m, err := handler(ctx, req)
		if !strings.Contains(info.FullMethod, "Ping") {
			global.Log.Infof("RPC: %s, client IP: %s, reply: %v, err: %v", info.FullMethod, ip, m, err)
		}
		return m, err
	}
}

// GetClientIP检查上下文以检索客户机的ip地址
func getClientIP(ctx context.Context) (string, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("couldn't parse client IP address")
	}
	return p.Addr.String(), nil
}
