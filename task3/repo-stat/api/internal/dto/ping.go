package dto

type ServiceStatusDTO struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponseDTO struct {
	Status   string             `json:"status"`
	Services []ServiceStatusDTO `json:"services"`
}
