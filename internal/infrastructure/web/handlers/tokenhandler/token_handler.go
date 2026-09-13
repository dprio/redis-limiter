package tokenhandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dprio/redis-limiter/internal/usecase/createtoken"
)

type TokenHandler struct {
	createTokenUseCase createtoken.UseCase
}

func New(createTokenUseCase createtoken.UseCase) *TokenHandler {
	return &TokenHandler{
		createTokenUseCase: createTokenUseCase,
	}
}

func (h *TokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.createTokenUseCase.Execute(ctx, request.ToCreateTokenInput())
	if errors.Is(err, createtoken.ErrInvalidTotalRequests) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := NewCreateTokenResponse(output)
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
