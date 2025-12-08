package dto

type SubscriptionDTO struct {
	ID          int    `json:"id"`
	ServiceName string `json:"service_name" validate:"required"`
	Price       int    `json:"price" validate:"gte=0"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	StartDate   string `json:"start_date" validate:"required,date_format"`
	EndDate     string `json:"end_date,omitempty" validate:"date_format"`
}
