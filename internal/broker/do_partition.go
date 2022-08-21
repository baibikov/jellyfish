// Package broker
/*
   Copyright 2022 Jellyfish message broker
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
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
