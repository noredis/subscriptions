package entity

import "time"

type Subscription struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

type SubscriptionFilter struct {
	Page        int
	Limit       int
	ServiceName string
	UserID      string
	StartDate   *time.Time
	EndDate     *time.Time
}
