package conn

import (
	"net"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Conn struct {
	net.Conn
}

func New(c net.Conn) *Conn {
	return &Conn{
		Conn: c,
	}
}

const maxMessageSize = 512

func (c Conn) WriteProto(m proto.Message) error {
	bb, err := proto.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "proto-marshal message")
	}

	_, err = c.Write(bb)
	return errors.Wrap(err, "write to connection")
}

func (c Conn) ReadProto(m proto.Message) error {
	bb := make([]byte, maxMessageSize)
	n, err := c.Read(bb)
	if err != nil {
		return errors.Wrap(err, "read from connection")
	}
	if n == 0 {
		return errors.New("undefined message from broker")
	}

	return errors.Wrap(proto.Unmarshal(bb[:n], m), "proto-unmarshal message")
}
