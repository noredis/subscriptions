package interfaces

import (
	"context"

	"github.com/noredis/subscriptions/internal/domain/entity"
)

type SubscriptionRepository interface {
	Insert(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
}
