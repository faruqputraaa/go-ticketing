package service

import (
	"context"
	"errors"
	_ "errors"
	"github.com/faruqputraaa/go-ticket/internal/entity"
	"github.com/faruqputraaa/go-ticket/internal/http/dto"
	"github.com/faruqputraaa/go-ticket/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type OfferService interface {
	GetAll(ctx context.Context) ([]entity.Offer, error)
	GetByID(ctx context.Context, id int64) (*entity.Offer, error)
	GetByIDUser(ctx context.Context, IDUser int) ([]entity.Offer, error)
	Create(ctx echo.Context, req dto.CreateOfferRequest) error
	Update(ctx context.Context, offer dto.UpdateOfferRequest) error
}

type offerService struct {
	offerRepository repository.OfferRepository
}

func NewOfferService(offerRepository repository.OfferRepository) OfferService {
	return &offerService{offerRepository}
}

func (s *offerService) Create(ctx echo.Context, req dto.CreateOfferRequest) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*entity.JWTCustomClaims)

	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.NameEvent == "" {
		return errors.New("name_event is required")
	}

	if req.Description == "" {
		return errors.New("description is required")
	}

	offer := &entity.Offer{
		IDUser:      claims.IDUser,
		Email:       req.Email,
		NameEvent:   req.NameEvent,
		Description: req.Description,
		Status:      "PENDING",
	}

	return s.offerRepository.Create(ctx.Request().Context(), offer)
}

// GetAll implements TicketService.
func (s *offerService) GetAll(ctx context.Context) ([]entity.Offer, error) {
	return s.offerRepository.GetAll(ctx)
}

// GetByID implements TicketService.
func (s *offerService) GetByID(ctx context.Context, id int64) (*entity.Offer, error) {
	return s.offerRepository.GetByID(ctx, id)
}

func (s *offerService) GetByIDUser(ctx context.Context, IDUser int) ([]entity.Offer, error) {
	return s.offerRepository.GetByIdUser(ctx, IDUser)
}

func (s *offerService) Update(ctx context.Context, req dto.UpdateOfferRequest) error {
	offer, err := s.offerRepository.GetByID(ctx, req.IDOffer)
	if err != nil {
		return err
	}

	if req.Status != "" {
		offer.Status = req.Status
	}

	return s.offerRepository.Update(ctx, offer)
}
