package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vanyaio/raketa-backend/internal/service"
	"github.com/vanyaio/raketa-backend/internal/storage"
	"github.com/vanyaio/raketa-backend/pkg/db"
	botpb "github.com/vanyaio/raketa-backend/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// db conn
	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// storage
	storage := storage.NewStorage(pool)

	// service
	service := service.NewBotService(storage)

	// grpc server
	go func() {
		if err := runGRPCServer(service); err != nil {
			log.Fatal(err)
		}
	}()

	// rest server
	go func() {
		if err := runRESTServer(service); err != nil {
			log.Fatal(err)
		}
	}()
	
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func runGRPCServer(service *service.BotService) error {
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))

	botpb.RegisterBotServiceServer(server, service)

	// reflection
	reflection.Register(server)
	// listen on port :50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("grpc server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	return nil
}

func runRESTServer(service *service.BotService) error {
	mux := runtime.NewServeMux()

	err := botpb.RegisterBotServiceHandlerServer(context.Background(), mux, service)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rest server start")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

	return nil
}
