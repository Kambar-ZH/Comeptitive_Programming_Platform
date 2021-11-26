package messagebroker

import "context"

type (
	MessageBroker interface {
		Connect(ctx context.Context) error
		Close() error

		Cache() CacheBroker
	}

	SubBrokerWithClient interface {
		Connect(ctx context.Context, brokers []string) error
		Close() error
	}
)
