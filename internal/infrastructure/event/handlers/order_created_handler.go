package handlers

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dprio/redis-limiter/internal/domain/eventtype"
	"github.com/dprio/redis-limiter/pkg/events"
	"github.com/streadway/amqp"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func CreateAndRegisterOrderCreatedHandler(rabbitMQChannel *amqp.Channel, dispatcher events.EventDispatcherInterface) events.EventHandlerInterface {
	handlr := &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}

	dispatcher.RegisterHandler(eventtype.OrderCreated, handlr)
	return handlr
}

func (h *OrderCreatedHandler) Handle(event events.Event, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
