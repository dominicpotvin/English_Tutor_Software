package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"englishtutor/internal/model"
)

// LevelProgress reports learner activity on lesson exercises, rolled up per
// level. Quiz questions are excluded; quiz performance is reported separately.
func (s *Store) LevelProgress(ctx context.Context) ([]model.LevelProgress, error) {
	rows, err := s.pool.Query(ctx, `
		select l.id as level_id, l.code, l.name,
		       count(distinct e.id)                                      as total_exercises,
		       count(distinct e.id) filter (where m.exercise_id is not null) as mastered_exercises,
		       count(a.id)                                               as total_attempts,
		       count(a.id) filter (where a.is_correct)                   as correct_attempts
		from levels l
		left join lessons les on les.level_id = l.id
		left join topics t    on t.lesson_id = les.id
		left join exercises e on e.topic_id = t.id
		left join attempts a  on a.exercise_id = e.id`+masteredJoin+`
		group by l.id
		order by l.position, l.id`)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.LevelProgress])
}
