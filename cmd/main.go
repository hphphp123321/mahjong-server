package main

import (
	"context"
	"github.com/hphphp123321/mahjong-server/app/bootloader"
	"github.com/hphphp123321/mahjong-server/app/bootloader/loaderlist/defaultloader"
	"github.com/hphphp123321/mahjong-server/app/dao"
	"github.com/hphphp123321/mahjong-server/app/global"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 加载模块，初始化global
	bootloader.Load(ctx, defaultloader.List)

	// graceful shutdown
	time.Sleep(time.Millisecond * 100)
	log.Println("------------------------------------------------")
	log.Printf("%s 已启动，按crtl+C停止程序.\n", global.C.App.Name)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-c:
		log.Println("Received signal to exit")
	}
	cancel()
	if err := dao.Close(); err != nil {
		log.Println("dao close error:", err)
	}
	time.Sleep(time.Second * 1)
}
