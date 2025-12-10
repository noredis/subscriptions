package dto

type CostFilterDTO struct {
	ServiceName string `json:"service_name"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date" validate:"required,date_format"`
	EndDate     string `json:"end_date" validate:"required,date_format"`
}

type TotalCostResponse struct {
	TotalCost int `json:"total_cost"`
}
