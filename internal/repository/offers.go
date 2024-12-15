package repository

import (
	"context"

	"github.com/faruqputraaa/go-ticket/internal/entity"
	"gorm.io/gorm"
)

type OfferRepository interface {
	GetAll(ctx context.Context) ([]entity.Offer, error)
	GetByID(ctx context.Context, id int64) (*entity.Offer, error)
	GetByIdUser(ctx context.Context, IDUser int64) ([]entity.Offer, error)

	Create(ctx context.Context, offer *entity.Offer) error
	Update(ctx context.Context, offer *entity.Offer) error
}

type offerRepository struct {
	db *gorm.DB
}

func NewOfferRepository(db *gorm.DB) OfferRepository {
	return &offerRepository{db}
}

func (r *offerRepository) GetAll(ctx context.Context) ([]entity.Offer, error) {
	result := make([]entity.Offer, 0)

	if err := r.db.WithContext(ctx).Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *offerRepository) GetByID(ctx context.Context, id int64) (*entity.Offer, error) {
	result := new(entity.Offer)
	if err := r.db.WithContext(ctx).Where("id_offer = ?", id).First(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *offerRepository) Create(ctx context.Context, offer *entity.Offer) error {
	return r.db.WithContext(ctx).Create(offer).Error
}

func (r *offerRepository) Update(ctx context.Context, offer *entity.Offer) error {
	return r.db.WithContext(ctx).
		Where("id_offer = ?", offer.IDOffer).
		Updates(offer).Error
}

func (r *offerRepository) GetByIdUser(ctx context.Context, IDUser int64) ([]entity.Offer, error) {
	var offers []entity.Offer
	err := r.db.WithContext(ctx).Where("id_user = ?", IDUser).Find(&offers).Error
	if err != nil {
		return nil, err
	}
	return offers, nil
}
