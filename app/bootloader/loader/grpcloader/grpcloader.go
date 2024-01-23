package grpcloader

import (
	"context"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/hphphp123321/mahjong-server/app/api/middleware/interceptor/authorization"
	"net"
	"strings"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1/mahjong"
	"github.com/hphphp123321/mahjong-server/app/controller"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// GrpcLoader 相关依赖：globals.ProjectRoot
// 1.初始化项目根目录globals.ProjectRoot
type GrpcLoader struct {
}

func (loader *GrpcLoader) Require() []string {
	return []string{"FINALLY"}
}

func (loader *GrpcLoader) Load(ctx context.Context, env map[string]string) error {
	var serverCfg = global.C.Server

	tcpAddr := fmt.Sprintf("%s:%d", serverCfg.IP, serverCfg.Port)
	global.Log.Debugln("Start listening at ", tcpAddr)
	lis, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(serverCfg.MaxConnectionIdle) * time.Second,
		MaxConnectionAgeGrace: time.Duration(serverCfg.MaxConnectionAgeGrace) * time.Second,
		Time:                  time.Duration(serverCfg.TimeTick) * time.Second,
		Timeout:               time.Duration(serverCfg.Timeout) * time.Second,
	}

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             time.Duration(serverCfg.MinTime) * time.Second,
		PermitWithoutStream: true,
	}

	// zap logger grpc中间件配置
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			if strings.Contains(fullMethodName, "Ping") {
				return false
			}
			return true
		},
		),
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			//logger.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(global.Log.Desugar(), zapOptions...),
			grpc_auth.UnaryServerInterceptor(authorization.AuthInterceptor),
		),
		grpc.ChainStreamInterceptor(
			//logger.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(global.Log.Desugar(), zapOptions...),
		),
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	)
	mahjongServer := controller.NewMahjongServer(server.NewImplServer())
	pb.RegisterMahjongServer(s, mahjongServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()
	return nil
}

func (loader *GrpcLoader) Name() string {
	return "GrpcLoader"
}

func StartRobotRegistryCheck() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			global.Log.Debugln("start check robot registry")
			server.TestAllRobots()
		}
	}()
}
