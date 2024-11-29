package sqlstore

import (
	"web/internal/app/model"
)

type TopicRepository struct {
	store *Store
}

func (r *TopicRepository) Create(topic *model.Topic) error {
	query := `
			INSERT INTO topics (user_id, name, description, content, is_public)
    	    VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at
	`

	return r.store.db.QueryRow(query, topic.UserID, topic.TopicName, topic.Description, topic.Content, topic.Visibility).
		Scan(&topic.ID, &topic.CreatedAt, &topic.UpdatedAt)
}

func (r *TopicRepository) FindByID(id int) (*model.Topic, error) {
	return nil, nil
}

func (r *TopicRepository) FindAll() ([]*model.Topic, error) { // аргменту заглушка чтобы удовлетворить интефейсу

	query := `SELECT id, user_id, name, description, content, is_public, created_at, updated_at
		FROM topics`

	rows, err := r.store.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []*model.Topic

	for rows.Next() {
		topic := &model.Topic{}

		err := rows.Scan(
			&topic.ID,
			&topic.UserID,
			&topic.TopicName,
			&topic.Description,
			&topic.Content,
			&topic.Visibility,
			&topic.CreatedAt,
			&topic.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}
