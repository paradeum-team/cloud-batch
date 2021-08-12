package main

import (
	"cloud-batch/api"
	"cloud-batch/configs"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// main query.collection.format默认为 csv , gin 框架 QueryArray 不支持， 修改为 multi
// @termsOfService https://gitlab.paradeum.com/pld/cloud-batch
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @query.collection.format multi
// @title Cloud Batch API
// @version v0.0.3
func main() {

	router := api.InitRouter()

	srv := &http.Server{
		Addr:           configs.Server.ListenAddr,
		Handler:        router,
		ReadTimeout:    time.Duration(configs.Server.DefaultReadTimeOut) * time.Second,
		WriteTimeout:   time.Duration(configs.Server.DefaultWriteTimeOut) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// service connections
		log.Printf("Actual pid is %d", syscall.Getpid())
		log.Printf("Now listening on: %s", configs.Server.ListenAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
