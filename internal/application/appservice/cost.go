package appservice

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/noredis/subscriptions/internal/application/dto"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/interfaces"
	"github.com/noredis/subscriptions/pkg/goext"
)

type CostService struct {
	validate *validator.Validate
	repo     interfaces.SubscriptionRepository
}

func NewCostService(
	validate *validator.Validate,
	repo interfaces.SubscriptionRepository,
) *CostService {
	return &CostService{
		validate: validate,
		repo:     repo,
	}
}

func (service *CostService) Total(
	ctx context.Context,
	f dto.CostFilterDTO,
) (int, error) {
	if err := service.validate.Struct(f); err != nil {
		return 0, err
	}

	filters, err := service.mapFiltersToEntity(f)
	if err != nil {
		return 0, err
	}

	subscriptions, err := service.repo.FindAll(ctx, filters)
	if err != nil {
		return 0, err
	}

	return service.calculateTotal(subscriptions, *filters.StartDate, *filters.EndDate), nil
}

func (service *CostService) calculateTotal(
	subs []*entity.Subscription,
	startDate time.Time,
	endDate time.Time,
) (total int) {
	for _, sub := range subs {
		total += service.calculate(sub, startDate, endDate)
	}
	return
}

func (service *CostService) calculate(
	sub *entity.Subscription,
	startDate time.Time,
	endDate time.Time,
) int {
	minDate := goext.MaxTime(startDate, sub.StartDate)
	maxDate := endDate

	if sub.EndDate != nil {
		maxDate = goext.MinTime(endDate, *sub.EndDate)
	}

	months := goext.MonthsBetween(minDate, maxDate)
	return months * sub.Price
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
