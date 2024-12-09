package service

import (
	"context"
	_ "errors"
	"fmt"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
)

type OfferService interface {
	GetAll(ctx context.Context) ([]entity.Offer, error)
	GetByID(ctx context.Context, id int64) (*entity.Offer, error)
	GetByIDUser(ctx context.Context, IDUser int64) ([]entity.Offer, error)
	Create(ctx context.Context, offer dto.CreateOfferRequest) error
	Update(ctx context.Context, offer dto.UpdateOfferRequest) error
}

type offerService struct {
	offerRepository repository.OfferRepository
}

func NewOfferService(offerRepository repository.OfferRepository) OfferService {
	return &offerService{offerRepository}
}

// Create implements TicketService.
func (s *offerService) Create(ctx context.Context, req dto.CreateOfferRequest) error {
	offer := &entity.Offer{
		IDUser:      req.IDUser,
		Email:       req.Email,
		NameEvent:   req.NameEvent,
		Description: req.Description,
		Status:      "PENDING",
	}

	fmt.Printf("Creating offer: %+v\n", offer)
	return s.offerRepository.Create(ctx, offer)
}

// GetAll implements TicketService.
func (s *offerService) GetAll(ctx context.Context) ([]entity.Offer, error) {
	return s.offerRepository.GetAll(ctx)
}

// GetByID implements TicketService.
func (s *offerService) GetByID(ctx context.Context, id int64) (*entity.Offer, error) {
	return s.offerRepository.GetByID(ctx, id)
}

func (s *offerService) GetByIDUser(ctx context.Context, IDUser int64) ([]entity.Offer, error) {
	return s.offerRepository.GetByIdUser(ctx, IDUser)
}

func (s *offerService) Update(ctx context.Context, req dto.UpdateOfferRequest) error {
	// Ambil offer berdasarkan ID
	offer, err := s.offerRepository.GetByID(ctx, req.IDOffer)
	if err != nil {
		return err
	}

	if req.Status != "" {
		offer.Status = req.Status
	}

	return s.offerRepository.Update(ctx, offer)
}
