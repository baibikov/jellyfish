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
	"time"

	"github.com/pkg/errors"

	"jellifish/internal/config"
)

type Listener struct {
	closed    bool
	listener  net.Listener
	broker    *Broker
	partition *Partition
	config    *config.Config
}

const tcpProtocol = "tcp"

func New(config *config.Config) (*Listener, error) {
	if config == nil {
		return nil, errors.New("config has not be empty")
	}

	l, err := net.Listen(tcpProtocol, config.Addr)
	if err != nil {
		return nil, errors.Wrap(err, "run listen broker")
	}

	listener := &Listener{
		listener: l,
		config:   config,
		broker:   NewBroker(),
	}

	if len(config.Slaves) != 0 {
		listener.partition, err = NewPartition(time.Second*2, config.Slaves)
		if err != nil {
			return nil, err
		}
	}

	return listener, err
}

func (l *Listener) Close() error {
	l.closed = true
	return l.listener.Close()
}

func (l *Listener) Broadcast(ctx context.Context) error {
	return l.broadcast(ctx)
}

func (l *Listener) broadcast(ctx context.Context) error {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			if l.closed {
				return nil
			}
			return errors.Wrap(err, "broadcast accept connection")
		}

		go NewHandler(conn, l.broker, l.partition).Do(ctx)
	}
}
