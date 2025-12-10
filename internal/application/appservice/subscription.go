package appservice

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/failure"
	"github.com/noredis/subscriptions/internal/domain/interfaces"
	"github.com/noredis/subscriptions/pkg/goext"
)

const dateFormat = "01-2006"

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
	req dto.SubscriptionRequest,
) (*dto.SubscriptionResponse, error) {
	if err := service.validate.Struct(req); err != nil {
		return nil, err
	}

	sub, err := service.mapToEntity(req)
	if err != nil {
		return nil, err
	}

	sub, err = service.repo.Insert(ctx, sub)
	if err != nil {
		return nil, err
	}

	return service.mapFromEntity(sub), nil
}

func (service *SubscriptionService) Update(
	ctx context.Context,
	req dto.SubscriptionRequest,
	id int,
) (*dto.SubscriptionResponse, error) {
	if err := service.validate.Struct(req); err != nil {
		return nil, err
	}

	sub, err := service.mapToEntity(req)
	if err != nil {
		return nil, err
	}

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
) (*dto.SubscriptionResponse, error) {
	sub, err := service.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return service.mapFromEntity(sub), nil
}

func (service *SubscriptionService) List(
	ctx context.Context,
	filters dto.SubscriptionFilterDTO,
) (*dto.SubscriptionListResponse, error) {
	f, err := service.mapFiltersToEntity(filters)
	if err != nil {
		return nil, err
	}

	if err := service.validate.Struct(f); err != nil {
		return nil, err
	}

	subscriptions, err := service.repo.Find(ctx, f)
	if err != nil {
		return nil, err
	}

	total, err := service.repo.Total(ctx, f)
	if err != nil {
		return nil, err
	}

	return &dto.SubscriptionListResponse{
		Page:  filters.Page,
		Limit: filters.Limit,
		Total: total,
		Data:  goext.Map(subscriptions, service.mapFromEntity),
	}, nil
}

func (service *SubscriptionService) mapToEntity(
	sub dto.SubscriptionRequest,
) (*entity.Subscription, error) {
	startDate, err := time.Parse(dateFormat, sub.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if sub.EndDate != "" {
		date, err := time.Parse(dateFormat, sub.EndDate)
		if err != nil {
			return nil, err
		}

		endDate = &date
	}

	return &entity.Subscription{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

func (service *SubscriptionService) mapFromEntity(
	sub *entity.Subscription,
) *dto.SubscriptionResponse {
	var endDate string
	if sub.EndDate != nil {
		endDate = sub.EndDate.Format(dateFormat)
	}

	return &dto.SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   sub.StartDate.Format(dateFormat),
		EndDate:     endDate,
	}
}

func (service *SubscriptionService) mapFiltersToEntity(
	f dto.SubscriptionFilterDTO,
) (*entity.SubscriptionFilter, error) {
	var startDate *time.Time
	if f.StartDate != "" {
		date, err := time.Parse(dateFormat, f.StartDate)
		if err != nil {
			return nil, err
		}

		startDate = &date
	}

	var endDate *time.Time
	if f.EndDate != "" {
		date, err := time.Parse(dateFormat, f.EndDate)
		if err != nil {
			return nil, err
		}

		endDate = &date
	}

	return &entity.SubscriptionFilter{
		Page:        f.Page,
		Limit:       f.Limit,
		ServiceName: f.ServiceName,
		UserID:      f.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
