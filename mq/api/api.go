package api

import "context"

type Message interface {
	Payload() []byte
	Topic() string
}

type Producer interface {
	Send(context.Context, []byte) error
}

type Consumer interface {
	Recv(context.Context) (Message, error)
	Ack(Message) error
}
type Client interface {
	Producer
	Consumer
}
