package apiservice

import (
	"encoding/json"
	"io"
)

func decodeResponse(response io.Reader) (ErrorResponse, error) {
	decoder := json.NewDecoder(response)
	var errorResponse ErrorResponse
	err := decoder.Decode(&errorResponse)
	if err != nil {
		return ErrorResponse{}, err
	}

	return errorResponse, nil
}

// ErrorResponse contains the errors details of a Google Analytics API response.
type ErrorResponse struct {
	Error struct {
		Errors []struct {
			Domain  string `json:"domain"`
			Reason  string `json:"reason"`
			Message string `json:"message"`
		} `json:"errors"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}
