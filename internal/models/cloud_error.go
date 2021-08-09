package models

type CloudError struct {
	Class   string `json:"class"`
	Code    int    `json:"code"`
	Details string `json:"details"`
}
