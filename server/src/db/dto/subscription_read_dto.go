package dto

type SubscriptionReadDTO struct {
	ServiceName *string `json:"service_name"`
	UserId      *string `json:"user_id"`
	Price       *uint   `json:"price"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}
