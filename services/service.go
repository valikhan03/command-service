package services

import (
	"context"

	"github.com/valikhan03/command-service/models"
	"github.com/valikhan03/command-service/pb"
	"github.com/valikhan03/command-service/repositories"

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
	id, err := s.repository.CreateAuction(ctx, req.Auction)
	if err != nil {
		return &pb.Response{}, err
	}
	req.Auction.Id = id
	event1 := models.Event{
		Command: tool.CRE_AUC_EVENT,
		Entity:  req.Auction,
	}
	s.eventChan <- &event1
	auctionFee, err := s.repository.FormAuctionFee(id)
	if err != nil{
		return &pb.Response{}, err
	}
	event2 := models.Event{
		Command: tool.CRE_INV_EVENT,
		Entity: auctionFee,
	}
	s.eventChan <- &event2
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
		Entity:  map[string]interface{}{"id":req.Id},
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) AddParticipant(ctx context.Context, req *pb.AddParticipantRequest) (*pb.Response, error) {
	requirements, err := s.repository.AddParticipant(ctx, req.AuctionId, req.ParticipantId)
	if err != nil {
		return &pb.Response{}, err
	}
	if requirements {
		return &pb.Response{
			Status: "",
		}, nil
	}

	attempt_reqs, err := s.repository.GetAttemptRequirements(ctx, req.AuctionId)
	if err != nil{
		return &pb.Response{}, err
	}

	if attempt_reqs.EnterFee>0 {
		event_inv := models.Event{
			Command: tool.CRE_INV_EVENT,
			Entity:   map[string]interface{}{"price": attempt_reqs.EnterFeeAmount, 
			"currency": "KZT",
			"user_id": req.ParticipantId, 
			"product_id": req.AuctionId,
			"product_type":2, 
			"auction_id": req.AuctionId},
		}
		s.eventChan <- &event_inv
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

func (s *Service) AddLot(ctx context.Context, req *pb.AddLotRequest) (*pb.Response, error) {
	id, err := s.repository.AddLot(ctx, req.Lot)
	if err != nil {
		return &pb.Response{}, err
	}

	req.Lot.Id = int32(id)
	event := models.Event{
		Command: tool.ADD_LOT_EVENT,
		Entity:  req.Lot,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) UpdateLot(ctx context.Context, req *pb.UpdateLotRequest) (*pb.Response, error) {
	err := s.repository.UpdateLot(ctx, req.Lot)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.UPD_LOT_EVENT,
		Entity:  req.Lot,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) DeleteLot(ctx context.Context, req *pb.DeleteLotRequest) (*pb.Response, error) {
	err := s.repository.DeleteLot(ctx, req.LotId)
	if err != nil {
		return &pb.Response{}, err
	}
	event := models.Event{
		Command: tool.DEL_LOT_EVENT,
		Entity:  req.LotId,
	}
	s.eventChan <- &event
	return &pb.Response{}, nil
}

func (s *Service) UploadMedia(stream pb.CommandService_UploadMediaServer) error {
	var mediaInfo *pb.MediaInfo

	for {
		req, err := stream.Recv()
		if err != nil{
			return err
		}

		mediaInfo = req.MediaInfo
	}

	err := s.repository.AddMediaInfo(context.TODO(), mediaInfo)
	if err != nil{
		return err
	}

	return nil
}