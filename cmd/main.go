package main

import (
	"net"
	"fmt"
	"os"
	"log"

	"auctions-service/db"
	"auctions-service/pb"
	"auctions-service/repositories"
	"auctions-service/services"

	"google.golang.org/grpc"
)

func main() {
	repository := repositories.NewRepository(db.InitDatabase())
	service := services.NewService(repository)

	grpcServer := grpc.NewServer()

	pb.RegisterAuctionsServiceServer(grpcServer, service)

	listener, err := net.Listen("tpc", fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
	if err != nil{
		log.Fatal(err)
	}

	err = grpcServer.Serve(listener)
	if err != nil{
		log.Fatal(err)
	}
}
