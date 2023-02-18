package services

import (
	"context"

	"auctions-service/models"
	"auctions-service/pb"
	"auctions-service/repositories"

	uuid "github.com/lithammer/shortuuid"
	"github.com/valikhan03/tool"
)

type Service struct {
	repository *repositories.Repository
	eventChan  chan<- *models.Event
}

func NewService(r *repositories.Repository, ch chan<- *models.Event) *Service {
	return &Service{
		repository: r,
		eventChan:  ch,
	}
}

func (s *Service) CreateAuction(ctx context.Context, req *pb.CreateAuctionRequest) (*pb.Response, error) {
	req.Auction.Id = uuid.New()
	err := s.repository.CreateAuction(ctx, req.Auction)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.CRE_AUC_EVENT,
		Entity:  req.Auction,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) UpdateAuction(ctx context.Context, req *pb.UpdateAuctionRequest) (*pb.Response, error) {
	err := s.repository.UpdateAuction(ctx, req.Auction)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.UPD_AUC_EVENT,
		Entity:  req.Auction,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) DeleteAuction(ctx context.Context, req *pb.DeleteAuctionRequest) (*pb.Response, error) {
	err := s.repository.DeleteAuction(ctx, req.Id)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.DEL_AUC_EVENT,
		Entity:  req.Id,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) AddParticipant(ctx context.Context, req *pb.AddParticipantRequest) (*pb.Response, error) {
	err := s.repository.AddParticipant(ctx, req.AuctionId, req.ParticipantId)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.ADD_PAR_EVENT,
		Entity:  map[string]interface{}{"auction_id": req.AuctionId, "participant_id": req.ParticipantId},
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) DeleteParticipant(ctx context.Context, req *pb.DeleteParticipantRequest) (*pb.Response, error) {
	err := s.repository.RemoveParticipant(ctx, req.AuctionId, req.ParticipantId)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.DEL_PAR_EVENT,
		Entity:  map[string]interface{}{"auction_id": req.AuctionId, "participant_id": req.ParticipantId},
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) AddProduct(ctx context.Context, req *pb.AddProductRequest) (*pb.Response, error) {
	err := s.repository.AddProduct(ctx, req.Product)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.ADD_LOT_EVENT,
		Entity:  req.Product,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Response, error) {
	err := s.repository.UpdateProduct(ctx, req.Product)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.UPD_LOT_EVENT,
		Entity:  req.Product,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) DeleteProduct(ctx context.Context, req *pb.DeletePoductRequest) (*pb.Response, error) {
	err := s.repository.DeleteProduct(ctx, req.ProductId)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.DEL_LOT_EVENT,
		Entity:  req.ProductId,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}