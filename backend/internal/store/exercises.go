package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"englishtutor/internal/model"
)

const exerciseCols = `id, topic_id, quiz_id, kind, prompt, choices, answer, explanation, position`

// ListExercisesByTopic returns the practice exercises of a topic, in order.
func (s *Store) ListExercisesByTopic(ctx context.Context, topicID int64) ([]model.Exercise, error) {
	rows, err := s.pool.Query(ctx,
		`select `+exerciseCols+` from exercises where topic_id = $1 order by position, id`, topicID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Exercise])
}

// ListExercisesByQuiz returns the questions of a quiz, in order.
func (s *Store) ListExercisesByQuiz(ctx context.Context, quizID int64) ([]model.Exercise, error) {
	rows, err := s.pool.Query(ctx,
		`select `+exerciseCols+` from exercises where quiz_id = $1 order by position, id`, quizID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Exercise])
}

// GetExercise returns a single exercise by id.
func (s *Store) GetExercise(ctx context.Context, id int64) (model.Exercise, error) {
	rows, err := s.pool.Query(ctx,
		`select `+exerciseCols+` from exercises where id = $1`, id)
	if err != nil {
		return model.Exercise{}, err
	}
	e, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Exercise])
	return e, notFound(err)
}

// CreateExercise inserts an exercise and returns the stored row.
func (s *Store) CreateExercise(ctx context.Context, in model.Exercise) (model.Exercise, error) {
	rows, err := s.pool.Query(ctx, `
		insert into exercises (topic_id, quiz_id, kind, prompt, choices, answer, explanation, position)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
		returning `+exerciseCols,
		in.TopicID, in.QuizID, in.Kind, in.Prompt, in.Choices, in.Answer, in.Explanation, in.Position)
	if err != nil {
		return model.Exercise{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Exercise])
}

// UpdateExercise overwrites an exercise and returns the stored row.
func (s *Store) UpdateExercise(ctx context.Context, id int64, in model.Exercise) (model.Exercise, error) {
	rows, err := s.pool.Query(ctx, `
		update exercises
		set topic_id = $2, quiz_id = $3, kind = $4, prompt = $5,
		    choices = $6, answer = $7, explanation = $8, position = $9
		where id = $1
		returning `+exerciseCols,
		id, in.TopicID, in.QuizID, in.Kind, in.Prompt, in.Choices, in.Answer, in.Explanation, in.Position)
	if err != nil {
		return model.Exercise{}, err
	}
	e, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Exercise])
	return e, notFound(err)
}

// DeleteExercise removes an exercise and its attempts (cascade).
func (s *Store) DeleteExercise(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `delete from exercises where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// SubmitAttempt grades the given answer against the exercise, records the
// attempt and returns both the attempt and the exercise it targeted.
func (s *Store) SubmitAttempt(ctx context.Context, exerciseID int64, given string) (model.Attempt, model.Exercise, error) {
	ex, err := s.GetExercise(ctx, exerciseID)
	if err != nil {
		return model.Attempt{}, model.Exercise{}, err
	}
	correct := model.Grade(ex, given)
	rows, err := s.pool.Query(ctx, `
		insert into attempts (exercise_id, given_answer, is_correct)
		values ($1, $2, $3)
		returning id, exercise_id, given_answer, is_correct, created_at`,
		exerciseID, given, correct)
	if err != nil {
		return model.Attempt{}, model.Exercise{}, err
	}
	att, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Attempt])
	if err != nil {
		return model.Attempt{}, model.Exercise{}, err
	}
	return att, ex, nil
}
