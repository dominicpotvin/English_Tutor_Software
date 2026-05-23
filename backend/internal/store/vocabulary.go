package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"englishtutor/internal/model"
)

const vocabularyCols = `id, level_id, lesson_id, category, term, definition, example, position`

// ListVocabulary returns vocabulary entries, optionally filtered by level,
// lesson or category. A nil filter argument means "no filter".
func (s *Store) ListVocabulary(ctx context.Context, levelID, lessonID *int64, category *string) ([]model.Vocabulary, error) {
	rows, err := s.pool.Query(ctx, `
		select `+vocabularyCols+`
		from vocabulary
		where ($1::bigint is null or level_id  = $1)
		  and ($2::bigint is null or lesson_id = $2)
		  and ($3::text   is null or category  = $3)
		order by category, position, id`, levelID, lessonID, category)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Vocabulary])
}

// CreateVocabulary inserts a vocabulary entry and returns the stored row.
func (s *Store) CreateVocabulary(ctx context.Context, in model.Vocabulary) (model.Vocabulary, error) {
	rows, err := s.pool.Query(ctx, `
		insert into vocabulary (level_id, lesson_id, category, term, definition, example, position)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning `+vocabularyCols,
		in.LevelID, in.LessonID, in.Category, in.Term, in.Definition, in.Example, in.Position)
	if err != nil {
		return model.Vocabulary{}, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Vocabulary])
}

// UpdateVocabulary overwrites a vocabulary entry and returns the stored row.
func (s *Store) UpdateVocabulary(ctx context.Context, id int64, in model.Vocabulary) (model.Vocabulary, error) {
	rows, err := s.pool.Query(ctx, `
		update vocabulary
		set level_id = $2, lesson_id = $3, category = $4, term = $5,
		    definition = $6, example = $7, position = $8
		where id = $1
		returning `+vocabularyCols,
		id, in.LevelID, in.LessonID, in.Category, in.Term, in.Definition, in.Example, in.Position)
	if err != nil {
		return model.Vocabulary{}, err
	}
	v, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Vocabulary])
	return v, notFound(err)
}

// DeleteVocabulary removes a vocabulary entry.
func (s *Store) DeleteVocabulary(ctx context.Context, id int64) error {
	tag, err := s.pool.Exec(ctx, `delete from vocabulary where id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
