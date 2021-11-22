package mq

import (
	"context"
	"errors"
	"github.com/sqjian/go-kit/mq/api"
	"github.com/sqjian/go-kit/mq/provider/pulsar"
	"strings"
)

type Instance struct {
	mqUrl  string
	client api.Client
}

func newDefaultInstanceConfig() *Instance {
	return &Instance{}
}

func NewMQInstance(mqType Type, mqUrl, mqTopic string, opts ...Option) (*Instance, error) {

	instance := newDefaultInstanceConfig()

	instance.mqUrl = mqUrl

	for _, opt := range opts {
		opt.apply(instance)
	}

	switch mqType {
	case Pulsar:
		{
			if !strings.HasPrefix(mqUrl, "pulsar") {
				return nil, ErrWrapper(ErrInvalidUrl)
			}
			if c, err := pulsar.NewPulsarClient(mqUrl, mqTopic, pulsar.WithDebug(true)); err != nil {
				return nil, err
			} else {
				instance.client = c
				return instance, nil
			}
		}
	default:
		{
			return nil, ErrWrapper(ErrInvalidUrl)
		}
	}
}

func (instance *Instance) DestroyInstance() {
	instance.client.Close()
}

func (instance *Instance) SendMsg(topic string, payload []byte) error {
	if _, ok := instance.client.; !ok {
		if producer, err := instance.client.NewProducer(topic); err == nil {
			instance.producers.Store(topic, producer)
		} else {
			return err
		}
	}

	val, _ := instance.producers.Load(topic)
	sender := val.(Producer)

	return sender.Send(context.Background(), payload)
}

func (instance *Instance) Ack(msg Message) error {
	topic := msg.Topic()
	if consumer, ok := instance.consumers.Load(topic); !ok {
		return errors.New(pulsar.ErrNoSuchConsumer)
	} else {
		return consumer.(Consumer).Ack(msg)
	}
}

func (instance *Instance) Nack(msg Message) error {
	topic := msg.Topic()
	if consumer, ok := instance.consumers.Load(topic); !ok {
		return errors.New(pulsar.ErrNoSuchConsumer)
	} else {
		return consumer.(Consumer).Nack(msg)
	}
}
