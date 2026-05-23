// Package model holds the data structures shared across the application:
// persisted entities, aggregated summaries and progress reports.
package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// StringList is a []string persisted as a JSON document (jsonb column).
type StringList []string

// Value serialises the list to a JSON string for storage in a jsonb column.
func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal([]string(s))
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan deserialises a JSON document from storage into the list.
func (s *StringList) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("model.StringList: cannot scan %T", src)
	}
	return json.Unmarshal(b, (*[]string)(s))
}

// --- Entities -------------------------------------------------------------

// Level is a curriculum track (Introduction to Grammar, Level 1, etc.).
type Level struct {
	ID          int64  `json:"id" db:"id"`
	Code        string `json:"code" db:"code"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Position    int    `json:"position" db:"position"`
}

// Lesson is an ordered unit inside a level.
type Lesson struct {
	ID       int64  `json:"id" db:"id"`
	LevelID  int64  `json:"levelId" db:"level_id"`
	Number   int    `json:"number" db:"number"`
	Title    string `json:"title" db:"title"`
	Summary  string `json:"summary" db:"summary"`
	Position int    `json:"position" db:"position"`
}

// Topic is a single teaching point inside a lesson; explanation is Markdown.
type Topic struct {
	ID          int64  `json:"id" db:"id"`
	LessonID    int64  `json:"lessonId" db:"lesson_id"`
	Title       string `json:"title" db:"title"`
	Explanation string `json:"explanation" db:"explanation"`
	Position    int    `json:"position" db:"position"`
}

// Quiz groups assessment exercises independently of the lesson tree.
type Quiz struct {
	ID          int64  `json:"id" db:"id"`
	LevelID     *int64 `json:"levelId" db:"level_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Position    int    `json:"position" db:"position"`
}

// Exercise is one practice item, owned by exactly one topic or one quiz.
type Exercise struct {
	ID          int64      `json:"id" db:"id"`
	TopicID     *int64     `json:"topicId" db:"topic_id"`
	QuizID      *int64     `json:"quizId" db:"quiz_id"`
	Kind        string     `json:"kind" db:"kind"`
	Prompt      string     `json:"prompt" db:"prompt"`
	Choices     StringList `json:"choices" db:"choices"`
	Answer      string     `json:"answer" db:"answer"`
	Explanation string     `json:"explanation" db:"explanation"`
	Position    int        `json:"position" db:"position"`
}

// Attempt records a single answer a learner submitted for an exercise.
type Attempt struct {
	ID          int64     `json:"id" db:"id"`
	ExerciseID  int64     `json:"exerciseId" db:"exercise_id"`
	GivenAnswer string    `json:"givenAnswer" db:"given_answer"`
	IsCorrect   bool      `json:"isCorrect" db:"is_correct"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// Vocabulary is a single word or expression to memorise.
type Vocabulary struct {
	ID         int64  `json:"id" db:"id"`
	LevelID    *int64 `json:"levelId" db:"level_id"`
	LessonID   *int64 `json:"lessonId" db:"lesson_id"`
	Category   string `json:"category" db:"category"`
	Term       string `json:"term" db:"term"`
	Definition string `json:"definition" db:"definition"`
	Example    string `json:"example" db:"example"`
	Position   int    `json:"position" db:"position"`
}

// --- Aggregates -----------------------------------------------------------

// LevelSummary is a level enriched with rollup counts for list views.
type LevelSummary struct {
	Level
	LessonCount       int `json:"lessonCount" db:"lesson_count"`
	TotalExercises    int `json:"totalExercises" db:"total_exercises"`
	MasteredExercises int `json:"masteredExercises" db:"mastered_exercises"`
}

// LessonSummary is a lesson enriched with rollup counts for list views.
type LessonSummary struct {
	Lesson
	TopicCount        int `json:"topicCount" db:"topic_count"`
	ExerciseCount     int `json:"exerciseCount" db:"exercise_count"`
	MasteredExercises int `json:"masteredExercises" db:"mastered_exercises"`
}

// TopicSummary is a topic enriched with rollup counts.
type TopicSummary struct {
	Topic
	ExerciseCount     int `json:"exerciseCount" db:"exercise_count"`
	MasteredExercises int `json:"masteredExercises" db:"mastered_exercises"`
}

// QuizSummary is a quiz enriched with question and mastery counts.
type QuizSummary struct {
	Quiz
	QuestionCount int `json:"questionCount" db:"question_count"`
	MasteredCount int `json:"masteredCount" db:"mastered_count"`
}

// LevelProgress reports learner activity rolled up to a level.
type LevelProgress struct {
	LevelID           int64  `json:"levelId" db:"level_id"`
	Code              string `json:"code" db:"code"`
	Name              string `json:"name" db:"name"`
	TotalExercises    int    `json:"totalExercises" db:"total_exercises"`
	MasteredExercises int    `json:"masteredExercises" db:"mastered_exercises"`
	TotalAttempts     int    `json:"totalAttempts" db:"total_attempts"`
	CorrectAttempts   int    `json:"correctAttempts" db:"correct_attempts"`
}
