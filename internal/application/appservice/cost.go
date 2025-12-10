package appservice

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/interfaces"
	"github.com/noredis/subscriptions/internal/domain/service"
)

type CostService struct {
	validate   *validator.Validate
	repo       interfaces.SubscriptionRepository
	calculator *service.CostCalculator
}

func NewCostService(
	validate *validator.Validate,
	repo interfaces.SubscriptionRepository,
	calculator *service.CostCalculator,
) *CostService {
	return &CostService{
		validate:   validate,
		repo:       repo,
		calculator: calculator,
	}
}

func (service *CostService) Total(
	ctx context.Context,
	f dto.CostFilterDTO,
) (*dto.TotalCostResponse, error) {
	if err := service.validate.Struct(f); err != nil {
		return nil, err
	}

	filters, err := service.mapFiltersToEntity(f)
	if err != nil {
		return nil, err
	}

	subscriptions, err := service.repo.FindAll(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &dto.TotalCostResponse{
		TotalCost: service.calculator.TotalCost(subscriptions, *filters.StartDate, *filters.EndDate),
	}, nil
}

func (service *CostService) mapFiltersToEntity(
	f dto.CostFilterDTO,
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
		ServiceName: f.ServiceName,
		UserID:      f.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
