package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"auctions-service/db"
	"auctions-service/kafka"
	"auctions-service/models"
	"auctions-service/pb"
	"auctions-service/repositories"
	"auctions-service/services"

	"google.golang.org/grpc"
)

func main() {
	eventsChan := make(chan *models.Event)

	repository := repositories.NewRepository(db.InitDatabase())
	service := services.NewService(repository, eventsChan)

	eventsProd := models.NewEventProducer(kafka.NewProducer(), eventsChan)
	go eventsProd.SendEvents()

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
