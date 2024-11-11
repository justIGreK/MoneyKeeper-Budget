package main

import (
	"context"
	"log"
	"net"

	"github.com/justIGreK/MoneyKeeper-Budget/cmd/handler"
	"github.com/justIGreK/MoneyKeeper-Budget/internal/repository"
	"github.com/justIGreK/MoneyKeeper-Budget/internal/service"
	"github.com/justIGreK/MoneyKeeper-Budget/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	user, err := client.NewUserClient("localhost:50052")
	if err != nil {
		log.Fatal(err)
	}
	db := repository.CreateMongoClient(ctx)
	budgetDB := repository.NewBudgetRepository(db)
	budgetSRV := service.NewBudgetService(budgetDB, user)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	handler := handler.NewHandler(grpcServer, budgetSRV)
	handler.RegisterServices()
	reflection.Register(grpcServer)

	log.Printf("Starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
