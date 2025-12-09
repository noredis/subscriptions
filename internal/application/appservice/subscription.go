package appservice

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/failure"
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

func (service *SubscriptionService) Create(
	ctx context.Context,
	req dto.SubscriptionDTO,
) (*dto.SubscriptionDTO, error) {
	if err := service.validate.Struct(req); err != nil {
		return nil, err
	}

	sub := service.mapToEntity(req)

	sub, err := service.repo.Insert(ctx, sub)
	if err != nil {
		return nil, err
	}

	return service.mapFromEntity(sub), nil
}

func (service *SubscriptionService) Update(
	ctx context.Context,
	req dto.SubscriptionDTO,
	id int,
) (*dto.SubscriptionDTO, error) {
	if err := service.validate.Struct(req); err != nil {
		return nil, err
	}

	sub := service.mapToEntity(req)

	exists, err := service.repo.ExistsByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, failure.ErrSubscriptionNotFound
	}

	sub.ID = id

	sub, err = service.repo.Update(ctx, sub)
	if err != nil {
		return nil, err
	}

	return service.mapFromEntity(sub), nil
}

func (service *SubscriptionService) Delete(ctx context.Context, id int) error {
	exists, err := service.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return failure.ErrSubscriptionNotFound
	}

	return service.repo.Delete(ctx, id)
}

func (service *SubscriptionService) Index(
	ctx context.Context,
	id int,
) (*dto.SubscriptionDTO, error) {
	sub, err := service.repo.FindByID(ctx, id)
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
