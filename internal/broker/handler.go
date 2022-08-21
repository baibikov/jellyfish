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
	"io"
	"net"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	pinger "jellyfish/pkg/ping"
	"jellyfish/protogenerated/messages"
)

type Handler struct {
	conn   net.Conn
	broker *Broker
	pp     *Partition
}

func NewHandler(conn net.Conn, broker *Broker, pp *Partition) *Handler {
	return &Handler{
		conn:   conn,
		broker: broker,
		pp:     pp,
	}
}

func (h *Handler) Do(ctx context.Context) {
	if err := h.do(ctx); err != nil {
		closeConnection(h.conn, "do")
		logrus.Error("broker", err)
	}
}

const (
	pingMessageSize    = 512
	produceMessageSize = 1024
)

func (h *Handler) do(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	buff := make([]byte, pingMessageSize)
	n, err := h.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "do connection read ping")
	}
	if n == 0 {
		return errNullMessage
	}

	ping := &messages.Ping{}
	err = proto.Unmarshal(buff[:n], ping)
	if err != nil {
		return errors.Wrap(err, "proto-message unmarshal")
	}

	pong := &messages.Pong{
		Pong: true,
	}

	bb, err := proto.Marshal(pong)
	if err != nil {
		return errors.Wrap(err, "proto-message marshal")
	}

	_, err = h.conn.Write(bb)
	if err != nil {
		return errors.Wrap(err, "do connection write pong")
	}

	logrus.Debugf("start listen messages by type %d", ping.GetPing())

	switch ping.Ping {
	case pinger.Publisher.Int32():
		logrus.Info("start producer execution")
		go h.producerDo(ctx)
	case pinger.Consumer.Int32():
		logrus.Info("start consumer execution")
		go h.consumerDo(ctx)
	case pinger.Partition.Int32():
		logrus.Info("start partition execution")
		go h.partitionDo(ctx)
	default:
		return errors.New("undefined ping message type")
	}

	return nil
}

func isSysError(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}

func closeConnection(conn net.Conn, space string) {
	if v := recover(); v != any(nil) {
		logrus.Errorf("space %s rec error: %+v", space, v)
	}

	err := conn.Close()
	if err != nil {
		logrus.Error(space, err)
	}
}

var errNullMessage = errors.New("null message")

func writeProto(conn net.Conn, message proto.Message) error {
	bb, err := proto.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "proto marshal")
	}

	_, err = conn.Write(bb)
	return errors.Wrap(err, "write proto message to connection")
}
