package api

import "context"

// UserRequest represents data provided by user.
type UserRequest struct {
	Value float64
	Curr  Currency
}

// ApiResponce represents responce from api.
type ApiResponce struct {
	Date     string  // date for current rate
	Value    float64 // actual rate
	CurrCode string  // currency code ex: 'USD'
	CurrName string  // full currenct name: ex: 'dolar amerykański'
}

// ApiProvider is common interfece for all apis.
type ApiProvider interface {
	GetCurrentRate(context.Context, UserRequest) (ApiResponce, error)
}
