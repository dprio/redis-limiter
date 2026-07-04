package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {

	t.Run("should create order successfully when all required parms ar passed", func(t *testing.T) {
		//given
		expectedOrder := Order{
			ID:    "batata",
			Price: 10.0,
			Tax:   1.0,
		}
		//when

		result, err := NewOrder(expectedOrder.ID, expectedOrder.Price, expectedOrder.Tax)

		//then
		assert.NoError(t, err)
		assert.Equal(t, expectedOrder, *result)
	})

	t.Run("should return error when invalid price", func(t *testing.T) {
		//given
		//when

		result, err := NewOrder("batata", -1, 10)

		//then
		assert.Nil(t, result)
		assert.EqualError(t, err, ErrInvalidPrice.Error())
	})

	t.Run("should return error when invalid tax", func(t *testing.T) {
		//given
		//when

		result, err := NewOrder("batata", 1, -10)

		//then
		assert.Nil(t, result)
		assert.EqualError(t, err, ErrInvalidTax.Error())
	})

	t.Run("should return error when invalid id", func(t *testing.T) {
		//given
		//when

		result, err := NewOrder("", 1, 10)

		//then
		assert.Nil(t, result)
		assert.EqualError(t, err, ErrInvalidID.Error())
	})
}

func TestCalculateFinalPrice(t *testing.T) {

	t.Run("should calculate final price successfully", func(t *testing.T) {
		//given
		expectedOrder := Order{
			ID:    "batata",
			Price: 10.0,
			Tax:   1.0,
		}

		expectedFinalPrice := expectedOrder.Price + expectedOrder.Tax
		//when

		err := expectedOrder.CaluculateFinalPeice()

		//then
		assert.NoError(t, err)
		assert.Equal(t, expectedFinalPrice, expectedOrder.FinalPrice)
	})

	t.Run("should return error when order is invalid", func(t *testing.T) {
		//given
		invalidOrder := Order{
			ID:    "batata",
			Price: -10.0,
			Tax:   1.0,
		}

		//when

		err := invalidOrder.CaluculateFinalPeice()

		//then
		assert.EqualError(t, err, ErrInvalidPrice.Error())
	})

}
