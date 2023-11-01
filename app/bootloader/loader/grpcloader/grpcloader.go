package grpcloader

import (
	"context"
	"fmt"
	"github.com/hphphp123321/mahjong-server/app/api/middleware/interceptor/logger"
	pb "github.com/hphphp123321/mahjong-server/app/api/v1/mahjong"
	"github.com/hphphp123321/mahjong-server/app/controller"
	"github.com/hphphp123321/mahjong-server/app/global"
	"github.com/hphphp123321/mahjong-server/app/service/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
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

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logger.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			logger.StreamServerInterceptor(),
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
