package eorder

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/m1ll3r1337/worker-service/internal/app/entity"
	ehandler "github.com/m1ll3r1337/worker-service/internal/app/handler/event"
	"github.com/m1ll3r1337/worker-service/internal/pkg/broker"
)

type handler struct{}

func NewHandler() ehandler.OrderCreated {
	return &handler{}
}

func (h *handler) CallbackOrderCreated(
	ctx context.Context, ev *entity.EventOrderCreated, _ []broker.Header,
) error {
	log.Info().Ctx(ctx).EmbedObject(ev).Msg("order.created event received")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(fmt.Errorf("order.created: empty order_guid"))
	}

	return nil
}
