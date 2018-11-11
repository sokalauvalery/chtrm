package dto

import "time"

type LoginRequest struct {
	Name             string `json:"name"`
	BaseURL          string `json:"baseURL"`
	HealthcheckRoute string `json:"healthcheckRoute"`
	MessageRoute     string `json:"messageRoute"`
}

type LoginResponse struct {
	UUID string `json:"uuid"`
}

type MessageRequest struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Time time.Time `json:"time"`
	Text string    `json:"text"`
}
