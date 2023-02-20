package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"github.com/valikhan03/command-service/db"
	"github.com/valikhan03/command-service/kafka"
	"github.com/valikhan03/command-service/models"
	"github.com/valikhan03/command-service/pb"
	"github.com/valikhan03/command-service/repositories"
	"github.com/valikhan03/command-service/services"
)

func main() {
	eventChan := make(chan *models.Event)
	repository := repositories.NewRepository(db.InitDatabase())
	service := services.NewService(repository, eventChan)

	producer := models.NewEventProducer(kafka.NewProducer(), eventChan)
	go producer.SendEvents()

	grpcServer := grpc.NewServer()

	pb.RegisterCommandServiceServer(grpcServer, service)

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Println(lis.Addr().String())

	grpcServer.Serve(lis)
}
