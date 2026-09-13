package tokenhandler

import "github.com/dprio/redis-limiter/internal/usecase/createtoken"

type createTokenRequest struct {
	TotalRequests int64 `json:"total_requests"`
}

func (r createTokenRequest) ToCreateTokenInput() createtoken.Input {
	return createtoken.Input{
		TotalRequests: r.TotalRequests,
	}
}

type createTokenResponse struct {
	Token         string `json:"token"`
	TotalRequests int64  `json:"total_requests"`
}

func NewCreateTokenResponse(output createtoken.Output) createTokenResponse {
	return createTokenResponse{
		Token:         output.Token,
		TotalRequests: output.TotalRequests,
	}
}
