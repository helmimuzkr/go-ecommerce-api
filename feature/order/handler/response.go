package handler

type paymentResponse struct {
	PaymentToken string `json:"token"`
	PaymentURL   string `json:"payment_link"`
}
