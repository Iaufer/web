package store

type Store interface {
	User() UserRepository
	Topic() TopicRepository
}
