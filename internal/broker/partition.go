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
	"go.uber.org/multierr"

	"jellyfish/internal/pkg/timeoutgroup"
	"jellyfish/pkg/conn"
	"jellyfish/pkg/ping"
	"jellyfish/protogenerated/messages"
)

// simple ISR message writing implement

type Partition struct {
	pullConnections []*peer
	timeout         time.Duration
}

func NewPartition(timeout time.Duration, connections []string) (*Partition, error) {
	if connections == nil {
		return nil, errors.New("broker connections has not be empty")
	}

	pull, err := initPullConnections(connections)
	if err != nil {
		return nil, err
	}

	return &Partition{
		pullConnections: pull,
		timeout:         timeout,
	}, nil
}

func initPullConnections(connections []string) (peers []*peer, err error) {
	pull := make([]*peer, len(connections))
	defer func() {
		if err == nil {
			return
		}

		for i, p := range pull {
			if p == nil {
				continue
			}

			multierr.AppendInto(&err, errors.Wrapf(p.Close(), "connection by index - [%d] close", i))
		}
	}()
	for i := 0; i < len(connections); i++ {
		if connections[i] == "" {
			return nil, errors.Errorf("connection by index - [%d] empty", i)
		}

		c, err := net.Dial(tcpProtocol, connections[i])
		if err != nil {
			return nil, errors.Wrapf(err, "connection by index - [%d]", i)
		}

		pull[i] = newPeer(c)
	}

	return pull, nil
}

type peer struct {
	isFail bool
	isInit bool

	net.Conn
}

func (p *peer) init() {
	p.isInit = true
}

func (p *peer) failed() {
	p.isFail = true
}

func (p *peer) success() {
	p.isFail = false
}

func newPeer(conn net.Conn) *peer {
	return &peer{
		Conn: conn,
	}
}

func (p *Partition) AskByPeers(ctx context.Context, m *messages.Partition) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	tg := timeoutgroup.New(ctx)

	tg.Go(func() error {
		for _, pp := range p.pullConnections {
			if pp.isFail {
				continue
			}

			if !pp.isInit && !pp.ping(ctx) {
				pp.failed()
				continue
			}

			if !pp.isInit {
				pp.init()
			}

			if err := pp.tryWrite(m); err != nil {
				pp.failed()
				return errors.Wrapf(err, "by remote addr - %s", pp.RemoteAddr())
			}

			pp.success()
		}
		return nil
	})

	return tg.Wait()
}

func (p *peer) ping(ctx context.Context) bool {
	return ping.New(p.Conn).Ping(ctx, ping.Partition) == nil
}

func (p *peer) tryWrite(m *messages.Partition) error {
	pc := conn.New(p)

	err := pc.WriteProto(m)
	if err != nil {
		return errors.Wrap(err, "partition write proto")
	}

	mm := &messages.PartitionAsk{}
	err = pc.ReadProto(mm)
	if err != nil {
		return errors.Wrap(err, "partition read proto")
	}
	if !mm.Ask {
		return errors.New("partition message not asked")
	}

	return nil
}
