package consumer

import (
	"encoding/json"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

type Payload struct {
	Message []byte
	err     error
}

func (p Payload) JsonUnmarshal(v interface{}) error {
	return errors.Wrap(
		json.Unmarshal(p.Message, v),
		"consumer payload json-unmarshal",
	)
}

func (p Payload) ProtoUnmarshal(m proto.Message) error {
	return errors.Wrap(
		proto.Unmarshal(p.Message, m),
		"consumer payload proto-unmarshal",
	)
}

func (p Payload) Err() error {
	return p.err
}
