package domain

import "errors"

var ErrTokenNotFound = errors.New("token not found")

type Token struct {
	ID            int64
	Hash          string
	TotalRequests int64
}
