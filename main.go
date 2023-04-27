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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/vanyaio/raketa-backend/internal/service"
	"github.com/vanyaio/raketa-backend/internal/storage"
	"github.com/vanyaio/raketa-backend/pkg/db"
	"github.com/vanyaio/raketa-backend/pkg/utils"
	proto "github.com/vanyaio/raketa-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = "GRPC_PORT"
	restPort = "REST_PORT"
)

func main() {
	grpcPort, err := utils.GetEnvValue(grpcPort)
	if err != nil {
		log.Fatal(err)
	}
	restPort, err := utils.GetEnvValue(restPort)
	if err != nil {
		log.Fatal(err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storage := storage.NewStorage(pool)

	service := service.NewBotService(storage)

	// grpc server
	go func() {
		if err := runGRPCServer(service, grpcPort); err != nil {
			log.Fatal(err)
		}
	}()

	// rest server
	go func() {
		if err := runRESTServer(service, restPort); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func runGRPCServer(service *service.Service, port string) error {
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))

	proto.RegisterRaketaServiceServer(server, service)

	reflection.Register(server)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("grpc server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	return nil
}

func runRESTServer(service *service.Service, port string) error {
	mux := runtime.NewServeMux()

	err := proto.RegisterRaketaServiceHandlerServer(context.Background(), mux, service)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("rest server start")
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}

	return nil
}
