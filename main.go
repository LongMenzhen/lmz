package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyrnicolase/lmz/config"
	_ "github.com/cyrnicolase/lmz/config"
	"github.com/cyrnicolase/lmz/engine"
	"github.com/cyrnicolase/lmz/route"
)

func main() {
	hub := engine.AttachHub()
	hub.Run()
	route.Route()

	srv := &http.Server{
		Addr: config.Config.HTTP.Addr,
	}
	go func() {
		if err := srv.ListenAndServe(); nil != err {
			log.Fatal("监听失败" + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 0*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); nil != err {
		log.Fatal("关闭服务器失败")
	}

	log.Println("服务器已关闭")
}
