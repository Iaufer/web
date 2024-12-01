package store

import "web/internal/app/model"

type TopicRepository interface {
	Create(topic *model.Topic) error
	FindByID(id int) (*model.Topic, error)
	FindAll() ([]*model.Topic, error)
	UpdateTopic(*model.Topic) error
}
