package createorder

import "github.com/dprio/redis-limiter/internal/domain"

type (
	Input struct {
		Price float64
		Tax   float64
	}

	Output struct {
		ID         string
		Price      float64
		Tax        float64
		FinalPrice float64
	}
)

func NewOutput(order domain.Order) Output {
	return Output{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}
}
