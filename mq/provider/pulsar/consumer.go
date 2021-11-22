package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/sqjian/go-kit/mq"
)

type Consumer struct {
	consumer pulsar.Consumer
}

func (pc *Consumer) Recv(ctx context.Context) (mq.Message, error) {
	msg, err := pc.consumer.Receive(ctx)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (pc *Consumer) Ack(msg mq.Message) error {
	pc.consumer.Ack(msg.(pulsar.Message))
	return nil
}

func (pc *Consumer) Nack(msg mq.Message) error {
	pc.consumer.Nack(msg.(pulsar.Message))
	return nil
}

func (pc *Consumer) Close() {
	pc.consumer.Close()
}
