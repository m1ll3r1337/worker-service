package eprocessor

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/m1ll3r1337/worker-service/internal/app/entity"
	ehandler "github.com/m1ll3r1337/worker-service/internal/app/handler/event"
	"github.com/m1ll3r1337/worker-service/internal/app/processor"
	"github.com/m1ll3r1337/worker-service/internal/pkg/broker"
)

type orderCreatedProc struct {
	h   ehandler.OrderCreated
	bus broker.Bus[entity.EventOrderCreated]
}

func NewOrderCreatedEventsCatcher(
	h ehandler.OrderCreated,
	bus broker.Bus[entity.EventOrderCreated],
) processor.Processor {
	return &orderCreatedProc{h, bus}
}

func (p *orderCreatedProc) StartAsync(ctx context.Context, wg *sync.WaitGroup) {
	err := p.bus.Subscribe(ctx, wg, p.h.CallbackOrderCreated)
	if err != nil {
		log.Fatal().Err(err).Str("topic_name", p.bus.QueueName()).Msg("failed to subscribe")
	}

	log.Info().Str("topic_name", p.bus.QueueName()).Msg("subscribed successfully")
}
