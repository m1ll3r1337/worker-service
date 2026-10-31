package eorder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/m1ll3r1337/worker-service/internal/app/entity"
	ehandler "github.com/m1ll3r1337/worker-service/internal/app/handler/event"
	"github.com/m1ll3r1337/worker-service/internal/app/service"
	"github.com/m1ll3r1337/worker-service/internal/pkg/broker"
)

type handler struct {
	delivery service.Delivery
	bus      broker.Bus[entity.EventOrderDeliveryCalculated]
}

func NewHandler(
	delivery service.Delivery,
	bus broker.Bus[entity.EventOrderDeliveryCalculated],
) ehandler.OrderCreated {
	return &handler{delivery: delivery, bus: bus}
}

func (h *handler) CallbackOrderCreated(
	ctx context.Context, ev *entity.EventOrderCreated, _ []broker.Header,
) error {
	log.Info().
		Ctx(ctx).
		EmbedObject(ev).
		Msg("Получено событие order.created")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(fmt.Errorf("order.created: empty order_guid"))
	}

	price, err := h.delivery.Calculate(ctx, ev.Currency)
	if err != nil {
		if errors.Is(err, entity.ErrFixerCurrencyNotFound) {
			return broker.NotCriticalError(fmt.Errorf("calculate delivery for %s: %w", ev.OrderGUID, err))
		}

		return fmt.Errorf("calculate delivery for %s: %w", ev.OrderGUID, err)
	}

	out := entity.EventOrderDeliveryCalculated{
		OrderGUID:     ev.OrderGUID,
		Currency:      ev.Currency,
		DeliveryPrice: price,
		CalculatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	err = h.bus.Send(ctx, &out,
		entity.BrokerHeaderOrderDeliveryCalculatedEventID(),
		entity.BrokerHeaderOrderDeliveryCalculatedType())
	if err != nil {
		return fmt.Errorf("publish delivery price for %s: %w", ev.OrderGUID, err)
	}

	log.Info().
		Ctx(ctx).
		EmbedObject(out).
		Msg("Опубликовано событие order.delivery.calculated")

	return nil
}
