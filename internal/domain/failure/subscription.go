package failure

import "errors"

var (
	ErrUserAlreadyHasThisSubscription = errors.New("user already has subscription to this service")
	ErrSubscriptionNotFound           = errors.New("subscription not found")
)
