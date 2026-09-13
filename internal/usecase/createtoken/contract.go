package createtoken

import "errors"

var ErrInvalidTotalRequests = errors.New("total-requests must be greater than zero")

type Input struct {
	TotalRequests int64
}

type Output struct {
	Token         string
	TotalRequests int64
}
