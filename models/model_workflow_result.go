package models

type WorkflowResult struct {
	Success       bool   `json:"Success,omitempty"`
	Status        int32  `json:"Status,omitempty"`
	Message       string `json:"Message,omitempty"`
	ActivityToken string `json:"ActivityToken,omitempty"`
}