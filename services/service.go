package services

import (
	"context"

	"auctions-service/pb"
	"auctions-service/repositories"
	"auctions-service/models"


	"github.com/google/uuid"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Service struct {
	repository *repositories.Repository
	eventChan chan <- *models.Event
}

func NewService(r *repositories.Repository, ch chan <- *models.Event) *Service {
	return &Service{
		repository: r,
		eventChan: ch,
	}
}

func (s *Service) CreateAuction(ctx context.Context, req *pb.CreateAuctionRequest) (*emptypb.Empty, error) {
	//get command from configs
	command := "CREATE_AUCTION"
	req.Auction.Id= uuid.New().String()
	err := s.repository.CreateAuction(ctx, req.Auction)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.Auction,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateAuction(ctx context.Context, req *pb.UpdateAuctionRequest) (*emptypb.Empty, error) {
	command := "UPDATE_AUCTION"
	err := s.repository.UpdateAuction(ctx, req.Auction)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.Auction,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil 
}
func (s *Service) CancelAuction(ctx context.Context, req *pb.CancelAuctionRequest) (*emptypb.Empty, error) {
	command := "CANCEL_AUCTION"
	err := s.repository.CancelAuction(ctx, req.Id)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.Id,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) AddParticipant(ctx context.Context, req *pb.AddParticipantRequest) (*emptypb.Empty, error) {
	command := "ADD_PARTICIPANT"
	err := s.repository.AddParticipant(ctx, req.AuctionId, req.ParticipantId)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: map[string]interface{}{"auction_id":req.AuctionId, "participant_id":req.ParticipantId},
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteParticipant(ctx context.Context, req *pb.DeleteParticipantRequest) (*emptypb.Empty, error) {
	command := "DELETE_PARTICIPANT"
	err := s.repository.RemoveParticipant(ctx, req.AuctionId, req.ParticipantId)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: map[string]interface{}{"auction_id":req.AuctionId, "participant_id":req.ParticipantId},
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*emptypb.Empty, error) {
	command := "ADD_PRODUCT"
	err := s.repository.AddProduct(ctx, req.Product)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.Product,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*emptypb.Empty, error) {
	command := "UPDATE_PRODUCT"
	err := s.repository.UpdateProduct(ctx, req.Product)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.Product,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteProduct(ctx context.Context, req *pb.DeletePoductRequest) (*emptypb.Empty, error) {
	command := "DELETE_PRODUCT"
	err := s.repository.DeleteProduct(ctx, req.ProductId)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.ProductId,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}