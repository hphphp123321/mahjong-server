package main

import (
	"flag"
	"fmt"
	pb "github.com/hphphp123321/mahjong-common/services/mahjong/v1"
	v1 "github.com/hphphp123321/mahjong-server/server"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

var (
	maxClients int
	address    string
	port       int

	minTime               int
	maxConnectionIdle     int
	maxConnectionAgeGrace int
	timeTick              int
	timeout               int

	logFormat string
	logLevel  string
	logOutput string
	logFile   string
)

func parseFlags() {
	flag.IntVar(&maxClients, "maxClients", 10, "max clients")
	flag.StringVar(&address, "address", "127.0.0.1", "server address")
	flag.IntVar(&port, "port", 16548, "port")

	flag.IntVar(&minTime, "minTime", 1, "If a client pings more than once every MinTime seconds, terminate the connection")
	flag.IntVar(&maxConnectionIdle, "maxConnectionIdle", 15, "If a client is idle for Idle seconds, send a GOAWAY")
	flag.IntVar(&maxConnectionAgeGrace, "maxConnectionAgeGrace", 5, "Allow Grace seconds for pending RPCs to complete before forcibly closing connections")
	flag.IntVar(&timeTick, "timeTick", 10, "Ping the client if it is idle for timeTick seconds to ensure the connection is still active")
	flag.IntVar(&timeout, "timeout", 5, "Wait 1 second for the ping ack before assuming the connection is dead")

	flag.StringVar(&logFormat, "logFormat", "text", "log format(json or text)")
	flag.StringVar(&logLevel, "logLevel", "debug", "log level(debug, info, warn, error, fatal, panic)")
	flag.StringVar(&logOutput, "logOutput", "stdout", "log output(stdout or stderr)")
	flag.StringVar(&logFile, "logFile", "", "log file path")
	flag.Parse()
}

func setupLogger() {
	switch logFormat {
	case "text":
		log.SetFormatter(&log.TextFormatter{
			ForceColors:               true,
			TimestampFormat:           "2006-01-02 15:04:05",
			FullTimestamp:             true,
			EnvironmentOverrideColors: true,
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				//处理文件名
				fileName := path.Base(frame.File)
				return ": " + strconv.Itoa(frame.Line), fileName
			},
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Error("set log format error")
	}

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.Error("set log level error")
	}

	switch logOutput {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.Error("set log output error")
	}

	if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
		} else {
			log.Error("set log file error")
		}
	}

}

func main() {
	parseFlags()
	setupLogger()
	log.Debug("Hello World!")
	tcpAddr := fmt.Sprintf("%s:%d", address, port)
	log.Debug("Start listening at ", tcpAddr)
	lis, err := net.Listen("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     time.Duration(maxConnectionIdle) * time.Second,
		MaxConnectionAgeGrace: time.Duration(maxConnectionAgeGrace) * time.Second,
		Time:                  time.Duration(timeTick) * time.Second,
		Timeout:               time.Duration(timeout) * time.Second,
	}

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             time.Duration(minTime) * time.Second,
		PermitWithoutStream: true,
	}
	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	server := v1.NewMahjongServer(10)
	pb.RegisterMahjongServer(s, server)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
