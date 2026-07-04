package orderhandler

import (
	"encoding/json"
	"net/http"

	"github.com/dprio/redis-limiter/internal/usecase/createorder"
	"github.com/dprio/redis-limiter/internal/usecase/getorders"
)

type OrderHandler struct {
	createOrderUseCase createorder.UseCase
	getOrderUseCase    getorders.UseCase
}

func New(createOrderUseCase createorder.UseCase, getOrdersUseCase getorders.UseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
		getOrderUseCase:    getOrdersUseCase,
	}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request createOrderRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, err := h.createOrderUseCase.Execute(ctx, request.ToCreateOrderInput())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := NewCreateOrderResponse(output)
	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *OrderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	output, err := h.getOrderUseCase.Execute(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]orderResponse, len(output))
	for i, out := range output {
		response[i] = NewCreateOrderResponse(createorder.Output(out))
	}

	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
