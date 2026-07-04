package createorder

import (
	"context"
	"testing"

	"github.com/dprio/redis-limiter/internal/domain"
	"github.com/dprio/redis-limiter/mocks/createorder"
	"github.com/dprio/redis-limiter/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecute(t *testing.T) {
	mockOrderRepository := createorder.NewMockOrderRepository(t)
	mockEventDispatcher := createorder.NewMockEventDispatcher(t)
	mockEventCreator := createorder.NewMockEventCreator(t)

	target := New(mockOrderRepository, mockEventDispatcher, mockEventCreator)

	cleanUp := func(t *testing.T) {
		mockOrderRepository.AssertExpectations(t)
		mockEventDispatcher.AssertExpectations(t)
		mockEventCreator.AssertExpectations(t)

		mockOrderRepository.Mock = mock.Mock{}
		mockEventDispatcher.Mock = mock.Mock{}
		mockEventCreator.Mock = mock.Mock{}
	}

	t.Run("should create order successfully", func(t *testing.T) {
		defer cleanUp(t)

		//given
		input := Input{
			Price: 10.0,
			Tax:   1.0,
		}

		eventCreator := events.NewEventCreator("")

		expectedOutput := Output{
			ID:         "",
			Price:      10.0,
			Tax:        1.0,
			FinalPrice: 11.0,
		}

		mockOrderRepository.EXPECT().
			Save(mock.Anything, mock.Anything).
			RunAndReturn(func(ctx context.Context, order *domain.Order) (domain.Order, error) {
				return *order, nil
			})

		mockEventCreator.EXPECT().Create(mock.Anything).Return(eventCreator.Create(expectedOutput))
		mockEventDispatcher.EXPECT().Dispatch(mock.Anything, mock.Anything).Return(nil)

		//when
		output, err := target.Execute(context.Background(), input)

		//then
		assert.NoError(t, err)
		assert.NotEmpty(t, output.ID)

		output.ID = ""
		assert.Equal(t, expectedOutput, output)
	})

	t.Run("should return error when invalid price value ", func(t *testing.T) {})

	t.Run("should return error when invalid tax value ", func(t *testing.T) {})

	t.Run("should return error when orderRepository returns error", func(t *testing.T) {})

	t.Run("should return error when eventDispatcher returns error", func(t *testing.T) {})
}
