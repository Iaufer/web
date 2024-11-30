package sqlstore

import (
	"database/sql"
	"web/internal/app/model"
	"web/internal/app/store"
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

	t := &model.Topic{}

	if err := r.store.db.QueryRow(
		"SELECT id, user_id, name, description, content, is_public, created_at, updated_at FROM topics WHERE id = $1",
		id,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.TopicName,
		&t.Description,
		&t.Content,
		&t.Visibility,
		&t.CreatedAt,
		&t.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return t, nil
}

func (r *TopicRepository) FindAll() ([]*model.Topic, error) {

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
