package pulsar

import (
	"context"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/sqjian/go-kit/mq"
	"regexp"
	"strings"
	"time"
)

func newDefaultClientConfig() *Client {
	cli := &Client{}
	cli.meta.debug = false
	cli.meta.MaxConnectionsPerBroker = 10
	cli.meta.ProducerWorkers = 3
	cli.meta.ConsumerWorkers = 3
	return cli
}

type ConsumerManager struct {
}
type ProducerManager struct {
}
type Client struct {
	meta struct {
		url   string
		topic string

		debug bool

		ProducerWorkers         int
		ConsumerWorkers         int
		MaxConnectionsPerBroker int
	}
	client pulsar.Client
}

func (pc *Client) Recv(ctx context.Context) (mq.Message, error) {
	panic("implement me")
}

func (pc *Client) Ack(message mq.Message) error {
	panic("implement me")
}

func (pc *Client) Send(ctx context.Context, bytes []byte) error {
	panic("implement me")
}

// "persistent://public/default/geo" -> geo
func (pc *Client) getRawTopicName(topic string) string {
	arr := strings.Split(topic, "/")
	return arr[len(arr)-1]
}

func (pc *Client) newProducer(topic string) (mq.Producer, error) {
	producer, err := pc.client.CreateProducer(pulsar.ProducerOptions{
		Topic:              topic,
		Name:               "prod-" + pc.getRawTopicName(topic),
		SendTimeout:        5 * time.Second,
		MaxPendingMessages: 10000,
		CompressionType:    pulsar.LZ4,
		CompressionLevel:   pulsar.Default,
	})
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, err
}

func (pc *Client) newConsumer(topic string) (mq.Consumer, error) {
	consumer, err := pc.client.Subscribe(pulsar.ConsumerOptions{
		Topic:                      topic,
		SubscriptionName:           "sub-" + pc.getRawTopicName(topic),
		Type:                       pulsar.Shared,
		ReplicateSubscriptionState: true,
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer}, nil
}

func NewPulsarClient(mqUrl string, topic string, opts ...Option) (*Client, error) {
	matched, err := regexp.Match(`^(persistent|non-persistent)://([a-z].*?/){2}[a-z].*$`, []byte(topic))
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, mq.ErrWrapper(mq.ErrTopicFormat)
	}

	cli := newDefaultClientConfig()
	cli.meta.url = mqUrl
	cli.meta.topic = topic

	for _, opt := range opts {
		opt.apply(cli)
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:                     cli.meta.url,
		MaxConnectionsPerBroker: cli.meta.MaxConnectionsPerBroker,
		Logger:                  logWrapper(cli.meta.debug),
	})
	if err != nil {
		return nil, err
	}
	cli.client = client

	return cli, nil
}

func (pc *Client) Close() {
	pc.client.Close()
}
