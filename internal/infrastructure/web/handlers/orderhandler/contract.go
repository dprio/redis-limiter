package orderhandler

import "github.com/dprio/redis-limiter/internal/usecase/createorder"

type (
	createOrderRequest struct {
		Price float64 `json:"price"`
		Tax   float64 `json:"tax"`
	}

	orderResponse struct {
		ID         string  `json:"id"`
		Price      float64 `json:"price"`
		Tax        float64 `json:"tax"`
		FinalPrice float64 `json:"final_price"`
	}
)

func (r *createOrderRequest) ToCreateOrderInput() createorder.Input {
	return createorder.Input{
		Price: r.Price,
		Tax:   r.Tax,
	}
}

func NewCreateOrderResponse(out createorder.Output) orderResponse {
	return orderResponse{
		ID:         out.ID,
		Price:      out.Price,
		Tax:        out.Tax,
		FinalPrice: out.FinalPrice,
	}
}
