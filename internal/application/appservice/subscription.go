package appservice

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/interfaces"
)

type SubscriptionService struct {
	validate *validator.Validate
	repo     interfaces.SubscriptionRepository
}

func NewSubscriptionService(
	validate *validator.Validate,
	repo interfaces.SubscriptionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		validate: validate,
		repo:     repo,
	}
}

func (service *SubscriptionService) CreateSubscription(
	ctx context.Context,
	req dto.SubscriptionDTO,
) (*dto.SubscriptionDTO, error) {
	if err := service.validate.Struct(req); err != nil {
		return nil, err
	}

	subscription := service.mapToEntity(req)

	sub, err := service.repo.Insert(ctx, subscription)
	if err != nil {
		return nil, err
	}

	return service.mapFromEntity(sub), nil
}

func (service *SubscriptionService) mapToEntity(
	sub dto.SubscriptionDTO,
) *entity.Subscription {
	return &entity.Subscription{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
}

func (service *SubscriptionService) mapFromEntity(
	sub *entity.Subscription,
) *dto.SubscriptionDTO {
	return &dto.SubscriptionDTO{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
	}
}
