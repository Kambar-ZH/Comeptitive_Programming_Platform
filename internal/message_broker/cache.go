package messagebroker

type CacheBroker interface {
	SubBrokerWithClient
	Remove(key interface{}) error
	Purge() error
}
