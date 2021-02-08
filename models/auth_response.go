package models

type AuthResponse struct {
	IdToken      string `json:"IdToken,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
}
