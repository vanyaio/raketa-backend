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
	"github.com/vanyaio/raketa-backend/config"
	"github.com/vanyaio/raketa-backend/internal/service"
	"github.com/vanyaio/raketa-backend/internal/storage"
	"github.com/vanyaio/raketa-backend/pkg/db"
	proto "github.com/vanyaio/raketa-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	filepath string = ".env"
)

func main() {
	config := config.GetConfig(".env")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	storage := storage.NewStorage(pool)

	service := service.NewBotService(storage)

	// grpc server
	go func() {
		if err := runGRPCServer(service, config); err != nil {
			log.Fatal(err)
		}
	}()

	// rest server
	go func() {
		if err := runRESTServer(service, config); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
}

func runGRPCServer(service *service.Service, config *config.Config) error {
	server := grpc.NewServer(grpc.MaxConcurrentStreams(1000))

	proto.RegisterRaketaServiceServer(server, service)

	reflection.Register(server)

	lis, err := net.Listen("tcp", config.GRPCServer.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("grpc server start")
	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
	return nil
}

func runRESTServer(service *service.Service, config *config.Config) error {
	gmux := runtime.NewServeMux()

	err := proto.RegisterRaketaServiceHandlerServer(context.Background(), gmux, service)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gmux)
	fs := http.FileServer(http.Dir("./swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	
	log.Println("rest server start")
	if err := http.ListenAndServe(config.RESTServer.RestPort, mux); err != nil {
		log.Fatal(err)
	}

	return nil
}
