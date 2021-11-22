package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
)

type Producer struct {
	producer pulsar.Producer
}

func (pp *Producer) Send(ctx context.Context, payload []byte) error {
	_, err := pp.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload:            payload,
		DisableReplication: false,
	})

	return err
}

func (pp *Producer) Close() {
	pp.producer.Close()
}
