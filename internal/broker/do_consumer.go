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
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"jellifish/protogenerated/messages"
)

func (h *Handler) consumerDo(ctx context.Context) {
	defer closeConnection(h.conn, "consumer")

	pp, err := consumerPayload(h.conn)
	if err != nil {
		if isSysError(err) {
			logrus.Info("consumer: close connection")
			return
		}

		logrus.Error("consumer: ", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := h.consumer(TopicName(pp.Topic)); err != nil {
				if isSysError(err) {
					logrus.Info("consumer: close connection")
					return
				}

				logrus.Error("consumer: ", err)
				return
			}
		}
	}
}

const (
	consumerMessageSize     = 1024
	consumerPingMessageSize = 1024
)

func consumerPayload(conn net.Conn) (*messages.ConsumerPayload, error) {
	buff := make([]byte, consumerMessageSize)
	n, err := conn.Read(buff)
	if err != nil {
		return nil, errors.Wrap(err, "read from buffer")
	}
	if n == 0 {
		return nil, errNullMessage
	}

	cp := &messages.ConsumerPayload{}
	err = proto.Unmarshal(buff[:n], cp)
	if err != nil {
		return nil, errors.Wrap(err, "proto-unmarshal consumer payload message")
	}

	_, err = conn.Write(buff[:n])
	if err != nil {
		return nil, errors.Wrap(err, "write to pong payload to connection")
	}

	return cp, nil
}

func (h *Handler) consumer(topic TopicName) error {
	buff := make([]byte, consumerPingMessageSize)
	n, err := h.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "read from buffer")
	}
	if n == 0 {
		return errNullMessage
	}

	mm := &messages.ConsumerResponse{}
	bb, err := h.broker.Read(topic)
	if err != nil {
		return errors.Wrapf(err, "read from broker by topic %s", topic)
	}
	if bb == nil {
		mm.IsEmpty = true
		return errors.Wrap(writeProto(h.conn, mm), "write zero message to connection")
	}

	mm.Message = bb
	return errors.Wrapf(
		writeProto(h.conn, mm),
		"write message by topic %s to connection",
		topic,
	)
}
