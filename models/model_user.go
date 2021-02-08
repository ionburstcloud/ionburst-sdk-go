package models

type User struct {
	Username     string `json:"Username,omitempty"`
	Password     string `json:"Password,omitempty"`
	RefreshToken string `json:"RefreshToken,omitempty"`
}
