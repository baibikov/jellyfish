package broker

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"jellifish/pkg/conn"
	"jellifish/protogenerated/messages"
)

func (h *Handler) partitionDo(ctx context.Context) {
	defer closeConnection(h.conn, "partition")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := h.partition(); err != nil {
				if isSysError(err) {
					logrus.Info("partition: close connection")
					return
				}

				logrus.Error("partition:", err)
				return
			}
		}
	}
}

func (h *Handler) partition() error {
	pc := conn.New(h.conn)

	mm := &messages.Partition{}
	err := pc.ReadProto(mm)
	if err != nil {
		return errors.Wrap(err, "read from partition")
	}

	err = pc.WriteProto(&messages.PartitionAsk{
		Ask: true,
	})
	if err != nil {
		return errors.Wrap(err, "write ask from partition")
	}

	logrus.Debugf("partition topic %s message %s asked", mm.Topic, mm.Message)
	return nil
}
