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
	"google.golang.org/protobuf/proto"

	"jellifish/protogenerated/messages"
)

func (h *Handler) producerDo(ctx context.Context) {
	defer closeConnection(h.conn, "producer")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := h.producer(ctx); err != nil {
				if isSysError(err) {
					logrus.Info("producer: close connection")
					return
				}

				logrus.Error("producer:", err)
				return
			}
		}
	}
}

func (h *Handler) producer(ctx context.Context) error {
	buff := make([]byte, produceMessageSize)
	n, err := h.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "read to buffer")
	}
	if n == 0 {
		return errNullMessage
	}

	pp := &messages.ProducerPayload{}
	err = proto.Unmarshal(buff[:n], pp)
	if err != nil {
		return errors.Wrap(err, "unmarshal message from buffer")
	}

	err = h.broker.Write(TopicName(pp.Topic), pp.Message)
	if err != nil {
		return errors.Wrap(err, "write message to broker")
	}

	ask := &messages.ProducerAsk{
		Ask: true,
	}

	if h.pp != nil {
		err = h.pp.AskByPeers(ctx, &messages.Partition{
			Topic:   pp.Topic,
			Message: pp.Message,
		})
		if err != nil {
			return err
		}
	}

	err = writeProto(h.conn, ask)
	if err != nil {
		return errors.Wrap(err, "ask message to connection")
	}

	logrus.Debugf("message %s by topic %s asked and saved", pp.Message, pp.Topic)
	return nil
}
