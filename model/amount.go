package model

type Amount struct {
	Value    string `json:"value" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}
