package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/Yershuaq/Asik_1_go/statistics/internal/usecase"
	events "github.com/Yershuaq/Asik_1_go/statistics/proto/events"
)

func Start(nc *nats.Conn, uc usecase.Usecase) error {
	for _, subj := range []string{"inventory.events", "order.events"} {
		if _, err := nc.Subscribe(subj, func(m *nats.Msg) {
			switch subj {
			case "inventory.events":
				var e events.InventoryEvent
				if err := proto.Unmarshal(m.Data, &e); err == nil {
					uc.HandleInventory(context.Background(), &e)
				}
			case "order.events":
				var e events.OrderEvent
				if err := proto.Unmarshal(m.Data, &e); err == nil {
					uc.HandleOrder(context.Background(), &e)
				}
			}
		}); err != nil {
			return err
		}
	}
	return nil
}
