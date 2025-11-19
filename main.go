package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/arthurhzna/Golang_gRPC/internal/handler"
	"github.com/arthurhzna/Golang_gRPC/pb/service"
	"github.com/arthurhzna/Golang_gRPC/pkg/database"
	"github.com/arthurhzna/Golang_gRPC/pkg/grpcmiddlerware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// db := database.ConnectDb(context.Background(), os.Getenv("DB_URL"))
	database.ConnectDb(context.Background(), os.Getenv("DB_URL"))

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		grpcmiddlerware.ErrorMiddleware,
	))

	if os.Getenv("ENVIRONMENT") == "DEV" {
		reflection.Register(grpcServer)
	}

	serviceHandler := handler.NewServiceHandler()

	service.RegisterHelloWorldServiceServer(grpcServer, serviceHandler)

	grpcServer.Serve(lis)

}
