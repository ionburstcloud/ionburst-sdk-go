package models

type Auth struct {
	Username     string `json:"Username"`
	Password     string `json:"Password,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
