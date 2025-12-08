package failure

import "errors"

var (
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
)
