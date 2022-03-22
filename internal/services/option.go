package services

import (
	messagebroker "site/internal/message_broker"
	"site/internal/store"

	lru "github.com/hashicorp/golang-lru"
)

type UserServiceOption func(u *UserServiceImpl)
type SubmissionServiceOption func(s *SubmissionServiceImpl)
type UploadFileServiceOption func(uf *UploadFileServiceImpl)
type AuthServiceOption func(a *AuthServiceImpl)
type ContestServiceOption func(c *ContestServiceImpl)
type ProblemServiceOption func(p *ProblemServiceImpl)
type ParticipantServiceOption func(p *ParticipantServiceImpl)
type UserFriendServiceOption func(u *UserFriendServiceImpl)

func UserServiceWithStore(store store.Store) UserServiceOption {
	return func(u *UserServiceImpl) {
		u.store = store
	}
}

func SubmissionServiceWithCache(cache *lru.TwoQueueCache) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.cache = cache
	}
}

func SubmissionServiceWithBroker(broker messagebroker.MessageBroker) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.broker = broker
	}
}

func SubmissionServiceWithStore(store store.Store) SubmissionServiceOption {
	return func(s *SubmissionServiceImpl) {
		s.store = store
	}
}

func UploadFileServiceWithStore(store store.Store) UploadFileServiceOption {
	return func(uf *UploadFileServiceImpl) {
		uf.store = store
	}
}

func AuthServiceWithStore(store store.Store) AuthServiceOption {
	return func(a *AuthServiceImpl) {
		a.store = store
	}
}

func ContestServiceWithStore(store store.Store) ContestServiceOption {
	return func(c *ContestServiceImpl) {
		c.store = store
	}
}

func ContestServiceWithBroker(broker messagebroker.MessageBroker) ContestServiceOption {
	return func(c *ContestServiceImpl) {
		c.broker = broker
	}
}

func ProblemServiceWithStore(store store.Store) ProblemServiceOption {
	return func(p *ProblemServiceImpl) {
		p.store = store
	}
}

func ParticipantServiceWithStore(store store.Store) ParticipantServiceOption {
	return func(p *ParticipantServiceImpl) {
		p.store = store
	}
}

func UserFriendsServiceWithStore(store store.Store) UserFriendServiceOption {
	return func(u *UserFriendServiceImpl) {
		u.store = store
	}
}
