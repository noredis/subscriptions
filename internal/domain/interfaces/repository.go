package interfaces

import (
	"context"

	"github.com/noredis/subscriptions/internal/domain/entity"
)

type SubscriptionRepository interface {
	Insert(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
	Update(ctx context.Context, subscription *entity.Subscription) (*entity.Subscription, error)
	Delete(ctx context.Context, id int) (error)
	ExistsByID(ctx context.Context, id int) (bool, error)
}
