package model

type Money struct {
	Amount int64  `json:"amount"`
	Code   string `json:"code"`
}

func NewMoney(amount int64, code string) *Money {
	return &Money{
		Amount: amount,
		Code:   code,
	}
}
