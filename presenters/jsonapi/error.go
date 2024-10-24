package jsonapi

import "encoding/json"

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
}

type ErrorResponse struct {
	Errors []*Error `json:"errors"`
}

func NewErrorResponse(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Errors: []*Error{{
			Title:  err.Error(),
			Status: status,
		}},
	}
}

func UnmarshalErrorResponse(bs []byte) (*ErrorResponse, error) {
	errResp := &ErrorResponse{}

	if err := json.Unmarshal(bs, errResp); err != nil {
		return nil, err
	}

	return errResp, nil
}
