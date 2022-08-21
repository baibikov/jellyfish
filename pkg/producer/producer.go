package producer

import (
	"context"
	"net"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"jellyfish/internal/pkg/timeoutgroup"
	"jellyfish/pkg/ping"
	"jellyfish/protogenerated/messages"
)

type Config struct {
	Addr string
}

type Producer struct {
	conn   net.Conn
	config *Config

	pinged bool
}

func New(config *Config) (*Producer, error) {
	if config == nil {
		return nil, errors.New("config has not empty")
	}

	conn, err := net.Dial("tcp", config.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "publisher: connect by addr %s", config.Addr)
	}

	return &Producer{
		conn:   conn,
		config: config,
	}, nil
}

func (p *Producer) Close() error {
	return errors.Wrap(p.conn.Close(), "publisher close")
}

type Params struct {
	Topic   string
	Message []byte
}

func (p *Producer) Push(ctx context.Context, params *Params) error {
	if params == nil {
		return nil
	}

	return p.push(ctx, params)
}

func (p *Producer) push(ctx context.Context, params *Params) error {
	if !p.pinged {
		if err := ping.New(p.conn).Ping(ctx, ping.Publisher); err != nil {
			return err
		}
		p.pinged = true
	}

	pp := &messages.ProducerPayload{
		Topic:   params.Topic,
		Message: params.Message,
	}

	bb, err := proto.Marshal(pp)
	if err != nil {
		return errors.Wrap(err, "proto-marshal producer payload")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	group := timeoutgroup.New(ctx)

	group.Go(func() error {
		return p.sendMessage(bb)
	})

	return errors.Wrap(group.Wait(), "message send")
}

const producerAskMessageSize = 512

func (p *Producer) sendMessage(payload []byte) error {
	_, err := p.conn.Write(payload)
	if err != nil {
		return errors.Wrap(err, "write payload to connection")
	}

	buff := make([]byte, producerAskMessageSize)
	n, err := p.conn.Read(buff)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("null from buffer reader")
	}

	ask := &messages.ProducerAsk{}
	err = proto.Unmarshal(buff[:n], ask)
	if err != nil {
		return errors.Wrap(err, "proto-unmarshal ask message")
	}
	if !ask.Ask {
		return errors.New("producer message dont asked")
	}

	return nil
}
