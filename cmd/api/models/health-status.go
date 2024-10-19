package models

type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}
