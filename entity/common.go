package entity

import (
	"encoding/json"
)

// ErrorResponse represents common error object
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents common success object
type SuccessResponse struct {
	Status string `json:"status"`
}

// NewError creates new object of ErrorResponse
func NewError(s string) ErrorResponse {
	return ErrorResponse{Error: s}
}

// NewErrorJSON creates new object of ErrorResponse and return byte array
func NewErrorJSON(s string) ([]byte, error) {
	e := ErrorResponse{Error: s}
	return json.Marshal(e)
}
