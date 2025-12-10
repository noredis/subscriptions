package service

import (
	"time"

	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/pkg/goext"
)

type CostCalculator struct{}

func NewCostCalculator() *CostCalculator {
	return &CostCalculator{}
}

func (calculator *CostCalculator) TotalCost(
	subs []*entity.Subscription,
	startDate time.Time,
	endDate time.Time,
) (total int) {
	for _, sub := range subs {
		total += calculator.SingleCost(sub, startDate, endDate)
	}
	return
}

func (calculator *CostCalculator) SingleCost(
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
