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

func NewService(r *repositories.Repository) *Service {
	return &Service{
		repository: r,
	}
}

func (s *Service) CreateAuction(ctx context.Context, req *pb.CreateAuctionRequest) (*emptypb.Empty, error) {
	//get command from configs
	command := "CREATE_AUCTION"
	req.Auction.ID= uuid.New().String()
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
	command := ""
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
	command := ""
	err := s.repository.CancelAuction(ctx, req.ID)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.ID,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) AddParticipant(ctx context.Context, req *pb.AddParticipantRequest) (*emptypb.Empty, error) {
	command := ""
	err := s.repository.AddParticipant(ctx, req.AuctionID, req.ParticipantID)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: map[string]interface{}{"auction_id":req.AuctionID, "participant_id":req.ParticipantID},
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) DeleteParticipant(ctx context.Context, req *pb.DeleteParticipantRequest) (*emptypb.Empty, error) {
	command := ""
	err := s.repository.RemoveParticipant(ctx, req.AuctionID, req.ParticipantID)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: map[string]interface{}{"auction_id":req.AuctionID, "participant_id":req.ParticipantID},
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}

func (s *Service) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*emptypb.Empty, error) {
	command := ""
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
	command := ""
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
	command := ""
	err := s.repository.DeleteProduct(ctx, req.ProductID)
	if err != nil{
		return &emptypb.Empty{}, err
	}
	event := models.Event{
		Command: command,
		Entity: req.ProductID,
	}
	s.eventChan <- &event
	return &emptypb.Empty{}, nil
}