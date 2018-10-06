package model

type ErrorResponse struct {
	Code int `json:"code"`
	Msg string `json:"message"`
}