package model

type ErrorResponse struct {
	Code int `json:"code"`
	Msg string `json:"message"`
}

type StringResponse struct {
	Code int `json:"code"`
	Msg string `json:"message"`
	Data string `json:"data"`
}