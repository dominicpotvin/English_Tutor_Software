// Package seed loads baseline curriculum content into the database.
//
// Two modes:
//
//   - First boot (no levels in DB): full seed — inserts every level, lesson,
//     topic, exercise, quiz and vocabulary entry from seed.json.
//   - Subsequent boots: fill mode — for each lesson in seed.json that has
//     topics defined, if the matching lesson in DB has zero topics, inserts
//     those topics and exercises. Existing content and learner progress are
//     left untouched. Quizzes and vocabulary are not touched in fill mode.
package seed

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed data/seed.json
var seedJSON []byte

type file struct {
	Levels     []levelSeed `json:"levels"`
	Quizzes    []quizSeed  `json:"quizzes"`
	Vocabulary []vocabSeed `json:"vocabulary"`
}

type levelSeed struct {
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Position    int          `json:"position"`
	Lessons     []lessonSeed `json:"lessons"`
}

type lessonSeed struct {
	Number   int         `json:"number"`
	Title    string      `json:"title"`
	Summary  string      `json:"summary"`
	Position int         `json:"position"`
	Topics   []topicSeed `json:"topics"`
}

type topicSeed struct {
	Title       string         `json:"title"`
	Explanation string         `json:"explanation"`
	Position    int            `json:"position"`
	Exercises   []exerciseSeed `json:"exercises"`
}

type exerciseSeed struct {
	Kind        string   `json:"kind"`
	Prompt      string   `json:"prompt"`
	Choices     []string `json:"choices"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
	Position    int      `json:"position"`
}

type quizSeed struct {
	LevelCode   string         `json:"level_code"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Position    int            `json:"position"`
	Questions   []exerciseSeed `json:"questions"`
}

type vocabSeed struct {
	LevelCode    string `json:"level_code"`
	LessonNumber *int   `json:"lesson_number"`
	Category     string `json:"category"`
	Term         string `json:"term"`
	Definition   string `json:"definition"`
	Example      string `json:"example"`
	Position     int    `json:"position"`
}

// Run loads seed content. On first boot it inserts everything; on later boots
// it only fills lessons that exist in DB but have no topics yet.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	var data file
	if err := json.Unmarshal(seedJSON, &data); err != nil {
		return fmt.Errorf("parse seed.json: %w", err)
	}
	if len(data.Levels) == 0 {
		return nil
	}

	var existing int
	if err := pool.QueryRow(ctx, `select count(*) from levels`).Scan(&existing); err != nil {
		return err
	}
	if existing == 0 {
		return runFullSeed(ctx, pool, &data)
	}
	return runFillEmptyLessons(ctx, pool, &data)
}

// runFullSeed inserts every level, lesson, topic, exercise, quiz and vocab
// entry from the seed file. Used only when the database is empty.
func runFullSeed(ctx context.Context, pool *pgxpool.Pool, data *file) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	levelID := map[string]int64{}
	lessonID := map[string]int64{}

	for _, lv := range data.Levels {
		var id int64
		if err := tx.QueryRow(ctx,
			`insert into levels (code, name, description, position) values ($1,$2,$3,$4) returning id`,
			lv.Code, lv.Name, lv.Description, lv.Position).Scan(&id); err != nil {
			return fmt.Errorf("insert level %s: %w", lv.Code, err)
		}
		levelID[lv.Code] = id

		for _, ls := range lv.Lessons {
			var lid int64
			if err := tx.QueryRow(ctx,
				`insert into lessons (level_id, number, title, summary, position) values ($1,$2,$3,$4,$5) returning id`,
				id, ls.Number, ls.Title, ls.Summary, ls.Position).Scan(&lid); err != nil {
				return fmt.Errorf("insert lesson %s/%d: %w", lv.Code, ls.Number, err)
			}
			lessonID[fmt.Sprintf("%s#%d", lv.Code, ls.Number)] = lid

			if err := insertLessonContent(ctx, tx, lid, ls.Topics); err != nil {
				return err
			}
		}
	}

	for _, qz := range data.Quizzes {
		var levelPtr *int64
		if id, ok := levelID[qz.LevelCode]; ok {
			levelPtr = &id
		}
		var qid int64
		if err := tx.QueryRow(ctx,
			`insert into quizzes (level_id, title, description, position) values ($1,$2,$3,$4) returning id`,
			levelPtr, qz.Title, qz.Description, qz.Position).Scan(&qid); err != nil {
			return fmt.Errorf("insert quiz %q: %w", qz.Title, err)
		}
		for _, q := range qz.Questions {
			if err := insertExercise(ctx, tx, nil, &qid, q); err != nil {
				return err
			}
		}
	}

	for _, vc := range data.Vocabulary {
		var levelPtr, lessonPtr *int64
		if id, ok := levelID[vc.LevelCode]; ok {
			levelPtr = &id
		}
		if vc.LessonNumber != nil {
			if id, ok := lessonID[fmt.Sprintf("%s#%d", vc.LevelCode, *vc.LessonNumber)]; ok {
				lessonPtr = &id
			}
		}
		if _, err := tx.Exec(ctx,
			`insert into vocabulary (level_id, lesson_id, category, term, definition, example, position)
			 values ($1,$2,$3,$4,$5,$6,$7)`,
			levelPtr, lessonPtr, vc.Category, vc.Term, vc.Definition, vc.Example, vc.Position); err != nil {
			return fmt.Errorf("insert vocabulary %q: %w", vc.Term, err)
		}
	}

	return tx.Commit(ctx)
}

// runFillEmptyLessons inserts topics and exercises for lessons that already
// exist in the database but have no topics yet. It never touches lessons that
// already have content, so learner progress is preserved.
func runFillEmptyLessons(ctx context.Context, pool *pgxpool.Pool, data *file) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	filled := 0
	for _, lv := range data.Levels {
		var levelID int64
		err := tx.QueryRow(ctx, `select id from levels where code = $1`, lv.Code).Scan(&levelID)
		if err == pgx.ErrNoRows {
			continue
		}
		if err != nil {
			return fmt.Errorf("lookup level %s: %w", lv.Code, err)
		}

		for _, ls := range lv.Lessons {
			if len(ls.Topics) == 0 {
				continue
			}
			var lessonID int64
			err := tx.QueryRow(ctx,
				`select id from lessons where level_id = $1 and number = $2`,
				levelID, ls.Number).Scan(&lessonID)
			if err == pgx.ErrNoRows {
				continue
			}
			if err != nil {
				return fmt.Errorf("lookup lesson %s/%d: %w", lv.Code, ls.Number, err)
			}

			var topicCount int
			if err := tx.QueryRow(ctx,
				`select count(*) from topics where lesson_id = $1`, lessonID).Scan(&topicCount); err != nil {
				return fmt.Errorf("count topics %s/%d: %w", lv.Code, ls.Number, err)
			}
			if topicCount > 0 {
				continue
			}

			if err := insertLessonContent(ctx, tx, lessonID, ls.Topics); err != nil {
				return err
			}
			filled++
		}
	}

	if filled == 0 {
		return tx.Rollback(ctx)
	}
	return tx.Commit(ctx)
}

func insertLessonContent(ctx context.Context, tx pgx.Tx, lessonID int64, topics []topicSeed) error {
	for _, tp := range topics {
		var tid int64
		if err := tx.QueryRow(ctx,
			`insert into topics (lesson_id, title, explanation, position) values ($1,$2,$3,$4) returning id`,
			lessonID, tp.Title, tp.Explanation, tp.Position).Scan(&tid); err != nil {
			return fmt.Errorf("insert topic %q: %w", tp.Title, err)
		}
		for _, ex := range tp.Exercises {
			if err := insertExercise(ctx, tx, &tid, nil, ex); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertExercise(ctx context.Context, tx pgx.Tx, topicID, quizID *int64, ex exerciseSeed) error {
	choices := ex.Choices
	if choices == nil {
		choices = []string{}
	}
	choicesJSON, err := json.Marshal(choices)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`insert into exercises (topic_id, quiz_id, kind, prompt, choices, answer, explanation, position)
		 values ($1,$2,$3,$4,$5,$6,$7,$8)`,
		topicID, quizID, ex.Kind, ex.Prompt, string(choicesJSON), ex.Answer, ex.Explanation, ex.Position)
	return err
}
