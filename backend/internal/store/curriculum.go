package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"englishtutor/internal/model"
)

// masteredJoin counts an exercise as mastered once it has any correct attempt.
const masteredJoin = `
	left join (select distinct exercise_id from attempts where is_correct) m
		on m.exercise_id = e.id`

// --- Levels ---------------------------------------------------------------

// ListLevelSummaries returns every level with rollup counts, ordered for display.
func (s *Store) ListLevelSummaries(ctx context.Context) ([]model.LevelSummary, error) {
	rows, err := s.pool.Query(ctx, `
		select l.id, l.code, l.name, l.description, l.position,
		       count(distinct les.id)                                    as lesson_count,
		       count(distinct e.id)                                      as total_exercises,
		       count(distinct e.id) filter (where m.exercise_id is not null) as mastered_exercises
		from levels l
		left join lessons les on les.level_id = l.id
		left join topics t    on t.lesson_id = les.id
		left join exercises e on e.topic_id = t.id`+masteredJoin+`
		group by l.id
		order by l.position, l.id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.LevelSummary])
}

// GetLevel returns a single level by id.
func (s *Store) GetLevel(ctx context.Context, id int64) (model.Level, error) {
	rows, err := s.pool.Query(ctx,
		`select id, code, name, description, position from levels where id = $1`, id)
	if err != nil {
		return model.Level{}, err
	}
	l, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Level])
	return l, notFound(err)
}

// CreateLevel inserts a level and returns the stored row.
func (s *Store) CreateLevel(ctx context.Context, in model.Level) (model.Level, error) {
	rows, err := s.pool.Query(ctx, `
		insert into levels (code, name, description, position)
		values ($1, $2, $3, $4)
		returning id, code, name, description, position`,
		in.Code, in.Name, in.Description, in.Position)
	if err != nil {
		return model.Level{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Level])
}

// UpdateLevel overwrites a level and returns the stored row.
func (s *Store) UpdateLevel(ctx context.Context, id int64, in model.Level) (model.Level, error) {
	rows, err := s.pool.Query(ctx, `
		update levels set code = $2, name = $3, description = $4, position = $5
		where id = $1
		returning id, code, name, description, position`,
		id, in.Code, in.Name, in.Description, in.Position)
	if err != nil {
		return model.Level{}, err
	}
	l, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Level])
	return l, notFound(err)
}

// --- Lessons --------------------------------------------------------------

// ListLessonSummaries returns the lessons of a level with rollup counts.
func (s *Store) ListLessonSummaries(ctx context.Context, levelID int64) ([]model.LessonSummary, error) {
	rows, err := s.pool.Query(ctx, `
		select les.id, les.level_id, les.number, les.title, les.summary, les.position,
		       count(distinct t.id)                                      as topic_count,
		       count(distinct e.id)                                      as exercise_count,
		       count(distinct e.id) filter (where m.exercise_id is not null) as mastered_exercises
		from lessons les
		left join topics t    on t.lesson_id = les.id
		left join exercises e on e.topic_id = t.id`+masteredJoin+`
		where les.level_id = $1
		group by les.id
		order by les.position, les.number, les.id`, levelID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.LessonSummary])
}

// GetLesson returns a single lesson by id.
func (s *Store) GetLesson(ctx context.Context, id int64) (model.Lesson, error) {
	rows, err := s.pool.Query(ctx,
		`select id, level_id, number, title, summary, position from lessons where id = $1`, id)
	if err != nil {
		return model.Lesson{}, err
	}
	l, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Lesson])
	return l, notFound(err)
}

// CreateLesson inserts a lesson and returns the stored row.
func (s *Store) CreateLesson(ctx context.Context, in model.Lesson) (model.Lesson, error) {
	rows, err := s.pool.Query(ctx, `
		insert into lessons (level_id, number, title, summary, position)
		values ($1, $2, $3, $4, $5)
		returning id, level_id, number, title, summary, position`,
		in.LevelID, in.Number, in.Title, in.Summary, in.Position)
	if err != nil {
		return model.Lesson{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Lesson])
}

// UpdateLesson overwrites a lesson and returns the stored row.
func (s *Store) UpdateLesson(ctx context.Context, id int64, in model.Lesson) (model.Lesson, error) {
	rows, err := s.pool.Query(ctx, `
		update lessons set level_id = $2, number = $3, title = $4, summary = $5, position = $6
		where id = $1
		returning id, level_id, number, title, summary, position`,
		id, in.LevelID, in.Number, in.Title, in.Summary, in.Position)
	if err != nil {
		return model.Lesson{}, err
	}
	l, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Lesson])
	return l, notFound(err)
}

// DeleteLesson removes a lesson and its topics and exercises (cascade).
func (s *Store) DeleteLesson(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `delete from lessons where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Topics ---------------------------------------------------------------

// ListTopicSummaries returns the topics of a lesson with rollup counts.
func (s *Store) ListTopicSummaries(ctx context.Context, lessonID int64) ([]model.TopicSummary, error) {
	rows, err := s.pool.Query(ctx, `
		select t.id, t.lesson_id, t.title, t.explanation, t.position,
		       count(distinct e.id)                                      as exercise_count,
		       count(distinct e.id) filter (where m.exercise_id is not null) as mastered_exercises
		from topics t
		left join exercises e on e.topic_id = t.id`+masteredJoin+`
		where t.lesson_id = $1
		group by t.id
		order by t.position, t.id`, lessonID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.TopicSummary])
}

// GetTopic returns a single topic by id.
func (s *Store) GetTopic(ctx context.Context, id int64) (model.Topic, error) {
	rows, err := s.pool.Query(ctx,
		`select id, lesson_id, title, explanation, position from topics where id = $1`, id)
	if err != nil {
		return model.Topic{}, err
	}
	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Topic])
	return t, notFound(err)
}

// CreateTopic inserts a topic and returns the stored row.
func (s *Store) CreateTopic(ctx context.Context, in model.Topic) (model.Topic, error) {
	rows, err := s.pool.Query(ctx, `
		insert into topics (lesson_id, title, explanation, position)
		values ($1, $2, $3, $4)
		returning id, lesson_id, title, explanation, position`,
		in.LessonID, in.Title, in.Explanation, in.Position)
	if err != nil {
		return model.Topic{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Topic])
}

// UpdateTopic overwrites a topic and returns the stored row.
func (s *Store) UpdateTopic(ctx context.Context, id int64, in model.Topic) (model.Topic, error) {
	rows, err := s.pool.Query(ctx, `
		update topics set lesson_id = $2, title = $3, explanation = $4, position = $5
		where id = $1
		returning id, lesson_id, title, explanation, position`,
		id, in.LessonID, in.Title, in.Explanation, in.Position)
	if err != nil {
		return model.Topic{}, err
	}
	t, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Topic])
	return t, notFound(err)
}

// DeleteTopic removes a topic and its exercises (cascade).
func (s *Store) DeleteTopic(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `delete from topics where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
