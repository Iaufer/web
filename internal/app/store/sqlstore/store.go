package sqlstore

import (
	"database/sql"
	"web/internal/app/store"

	_ "github.com/lib/pq" //...
)

type Store struct {
	db              *sql.DB
	userRepository  *UserRepository
	topicRepository *TopicRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Topic() store.TopicRepository {
	if s.topicRepository == nil {
		s.topicRepository = &TopicRepository{store: s} // Инициализируем TopicRepository
	}
	return s.topicRepository
}
