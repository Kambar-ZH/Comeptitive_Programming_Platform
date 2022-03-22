package messagebroker

import "site/internal/datastruct"

type ContestBroker interface {
	SubBrokerWithClient
	CreateContest(contest *datastruct.Contest) error
}
