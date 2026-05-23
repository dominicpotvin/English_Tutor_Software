package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"englishtutor/internal/model"
)

// ListQuizSummaries returns every quiz with question and mastery counts.
func (s *Store) ListQuizSummaries(ctx context.Context) ([]model.QuizSummary, error) {
	rows, err := s.pool.Query(ctx, `
		select q.id, q.level_id, q.title, q.description, q.position,
		       count(e.id)                                               as question_count,
		       count(distinct e.id) filter (where m.exercise_id is not null) as mastered_count
		from quizzes q
		left join exercises e on e.quiz_id = q.id`+masteredJoin+`
		group by q.id
		order by q.position, q.id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.QuizSummary])
}

// GetQuiz returns a single quiz by id.
func (s *Store) GetQuiz(ctx context.Context, id int64) (model.Quiz, error) {
	rows, err := s.pool.Query(ctx,
		`select id, level_id, title, description, position from quizzes where id = $1`, id)
	if err != nil {
		return model.Quiz{}, err
	}
	q, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Quiz])
	return q, notFound(err)
}

// CreateQuiz inserts a quiz and returns the stored row.
func (s *Store) CreateQuiz(ctx context.Context, in model.Quiz) (model.Quiz, error) {
	rows, err := s.pool.Query(ctx, `
		insert into quizzes (level_id, title, description, position)
		values ($1, $2, $3, $4)
		returning id, level_id, title, description, position`,
		in.LevelID, in.Title, in.Description, in.Position)
	if err != nil {
		return model.Quiz{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Quiz])
}

// UpdateQuiz overwrites a quiz and returns the stored row.
func (s *Store) UpdateQuiz(ctx context.Context, id int64, in model.Quiz) (model.Quiz, error) {
	rows, err := s.pool.Query(ctx, `
		update quizzes set level_id = $2, title = $3, description = $4, position = $5
		where id = $1
		returning id, level_id, title, description, position`,
		id, in.LevelID, in.Title, in.Description, in.Position)
	if err != nil {
		return model.Quiz{}, err
	}
	q, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Quiz])
	return q, notFound(err)
}

// DeleteQuiz removes a quiz and its questions (cascade).
func (s *Store) DeleteQuiz(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `delete from quizzes where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
