package consumer

import (
	"context"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/baibikov/jellyfish/pkg/conn"
	"github.com/baibikov/jellyfish/pkg/ping"
	"github.com/baibikov/jellyfish/protogenerated/messages"
)

type Consumer struct {
	conn   *conn.Conn
	config *Config
	closed bool
	pinged bool

	payload chan Payload
	once    sync.Once
}

func (c *Consumer) writeMessage(b []byte) {
	if c.closed {
		return
	}

	c.payload <- Payload{
		Message: b,
	}
}

func (c *Consumer) writeError(err error) {
	if c.closed {
		return
	}

	c.payload <- Payload{
		err: err,
	}
}

type Config struct {
	Addr string
}

func New(config *Config) (*Consumer, error) {
	if config == nil {
		return nil, errors.New("config has not empty")
	}

	nc, err := net.Dial("tcp", config.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "consumer: connect by addr %s", config.Addr)
	}

	c := &Consumer{
		conn:    conn.New(nc),
		config:  config,
		payload: make(chan Payload),
	}

	return c, nil
}

func (c *Consumer) Close() error {
	c.once.Do(func() {
		close(c.payload)
		c.closed = true
	})
	return errors.Wrap(c.conn.Close(), "consumer close")
}

func (c *Consumer) Consume(ctx context.Context, topic string) <-chan Payload {
	go c.do(ctx, topic)
	return c.payload
}

func (c *Consumer) do(ctx context.Context, topic string) {
	if err := c.broadcast(ctx, topic); err != nil {
		c.writeError(err)

		logrus.Error("message from broadcast: ", err)
	}
}

func (c *Consumer) broadcast(ctx context.Context, topic string) error {
	if !c.pinged {
		if err := ping.New(c.conn).Ping(ctx, ping.Consumer); err != nil {
			return err
		}
		c.pinged = true
	}

	err := c.conn.WriteProto(&messages.ConsumerPayload{Topic: topic})
	if err != nil {
		return errors.Wrap(err, "proto-marshal message")
	}

	dest := make([]byte, 512)
	_, err = c.conn.Read(dest)
	if err != nil {
		return errors.Wrap(err, "read from broker")
	}

	message := &messages.ConsumerResponse{}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err = c.conn.Write([]byte("1"))
			if err != nil {
				return errors.Wrap(err, "ping consumer broker")
			}

			err = c.conn.ReadProto(message)
			if err != nil {
				return errors.Wrap(err, "from broker")
			}

			if message.IsEmpty {
				continue
			}

			c.writeMessage(message.GetMessage())
			message.Reset()
		}
	}
}
