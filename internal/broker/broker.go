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
	"sync"

	"github.com/pkg/errors"
)

type Broker struct {
	mutex sync.RWMutex
	topic *Topic
}

func NewBroker() *Broker {
	return &Broker{
		topic: &Topic{
			mp: make(map[TopicName]*pack),
		},
	}
}

func (b *Broker) Write(name TopicName, payload []byte) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if !b.topic.exists(name) {
		if err := b.topic.create(name); err != nil {
			return err
		}
		return nil
	}

	b.topic.pack(name).append(payload)
	return nil
}

func (b *Broker) Read(name TopicName) ([]byte, error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if !b.topic.exists(name) {
		if err := b.topic.create(name); err != nil {
			return nil, err
		}
	}

	return b.topic.pack(name).message(), nil
}

type Topic struct {
	mp map[TopicName]*pack
}

func (t Topic) exists(name TopicName) bool {
	if t.mp == nil {
		return false
	}

	_, ok := t.mp[name]
	return ok
}

func (t *Topic) create(name TopicName) error {
	if t.mp == nil {
		return errors.New("topic storage not initialized")
	}

	t.mp[name] = &pack{}
	return nil
}

func (t *Topic) pack(name TopicName) *pack {
	vv, _ := t.mp[name]
	if vv == nil {
		p := &pack{}
		t.mp[name] = p
		vv = p
	}

	return vv
}

type pack struct {
	messages    [][]byte
	writeOffset int
	readOffset  int
}

func (m *pack) message() []byte {
	if len(m.messages) == 0 {
		return nil
	}

	if m.readOffset > len(m.messages)-1 {
		return nil
	}

	p := m.messages[m.readOffset]
	m.readOffset++
	return p
}

func (m *pack) append(payload []byte) {
	m.messages = append(m.messages, payload)
	m.writeOffset++
}

type TopicName string
