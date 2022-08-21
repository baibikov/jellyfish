package ping

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/baibikov/jellyfish/protogenerated/messages"
)

type Ping struct {
	conn net.Conn
}

func New(conn net.Conn) Pinger {
	return &Ping{
		conn: conn,
	}
}

type Pinger interface {
	Ping(ctx context.Context, pt PayloadType) error
}

type PayloadType int

func (p PayloadType) Int32() int32 {
	return int32(p)
}

const (
	Publisher PayloadType = iota + 1
	Consumer
	Partition
)

func (p *Ping) Ping(ctx context.Context, pt PayloadType) error {
	select {
	case <-ctx.Done():
		return nil
	default:

	}

	return p.ping(pt)
}

func (p *Ping) ping(pt PayloadType) error {
	message := messages.Ping{
		Ping: pt.Int32(),
	}
	bb, err := proto.Marshal(&message)
	if err != nil {
		return errors.Wrap(err, "proto-marsh ping message")
	}
	if len(bb) == 0 {
		return errors.New("proto-marshal null value")
	}

	_, err = p.conn.Write(bb)
	if err != nil {
		return errors.Wrap(err, "write proto-message")
	}

	buff := make([]byte, 512)
	n, err := p.conn.Read(buff)
	if err != nil {
		return errors.Wrap(err, "read pong message")
	}
	if n == 0 {
		return errors.New("null from buffer reader")
	}

	pong := &messages.Pong{}
	err = proto.Unmarshal(buff[:n], pong)
	if err != nil {
		return errors.Wrap(err, "unmarshal pong proto-message")
	}
	if !pong.Pong {
		return errors.New("ping message not ponged")
	}

	return nil
}
