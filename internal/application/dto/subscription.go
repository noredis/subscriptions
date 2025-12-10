package dto

type SubscriptionRequest struct {
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"gte=0"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required,date_format"`
	EndDate     string `json:"end_date,omitempty" validate:"date_format"`
}

type SubscriptionResponse struct {
	ID          int    `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date,omitempty"`
}

type SubscriptionListResponse struct {
	Page  int                     `json:"page"`
	Limit int                     `json:"limit"`
	Total int                     `json:"total"`
	Data  []*SubscriptionResponse `json:"data"`
}

type SubscriptionFilterDTO struct {
	Page        int `validate:"gte=1"`
	Limit       int `validate:"gte=1"`
	ServiceName string
	UserID      string
	StartDate   string `validate:"date_format"`
	EndDate     string `validate:"date_format"`
}
