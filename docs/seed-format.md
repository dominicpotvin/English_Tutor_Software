# Seed content format

The backend seeds baseline curriculum content from `backend/internal/seed/data/seed.json`
(embedded into the binary).

## Seeder modes

The seeder runs at every backend start with two possible modes:

- **First boot** — `levels` table is empty. The seeder performs a **full insert**: every
  level, lesson, topic, exercise, quiz and vocabulary entry from `seed.json` is loaded.
- **Subsequent boots** — `levels` table is not empty. The seeder switches to **fill mode**:
  for each lesson defined in `seed.json` that has `topics`, if the matching lesson in DB
  has zero topics, those topics and exercises are inserted. Lessons that already have
  content are not touched, so learner progress is preserved. Quizzes and vocabulary are
  not modified in this mode.

Practical consequences:

- Adding topics to a previously-empty lesson in `seed.json` and restarting the backend
  is enough — the new topics are loaded automatically without resetting the database.
- Editing existing topics in `seed.json` does **not** update them in DB. Use the API or
  the MCP connector for content updates, or run `docker compose down -v` to wipe and
  reseed.

## Structure

```json
{
  "levels": [
    {
      "code": "IG",
      "name": "Introduction to Grammar",
      "description": "Short level description.",
      "position": 1,
      "lessons": [
        {
          "number": 1,
          "title": "Lesson title",
          "summary": "One or two sentences.",
          "position": 1,
          "topics": [
            {
              "title": "Topic title",
              "explanation": "Markdown teaching content.",
              "position": 1,
              "exercises": [
                {
                  "kind": "mcq",
                  "prompt": "A ___ is playing with a ball.",
                  "choices": ["girl", "girls", "girles"],
                  "answer": "girl",
                  "explanation": "Why this answer is correct.",
                  "position": 1
                }
              ]
            }
          ]
        }
      ]
    }
  ],
  "quizzes": [
    {
      "level_code": "IG",
      "title": "Quiz title",
      "description": "Short description.",
      "position": 1,
      "questions": [
        {
          "kind": "mcq",
          "prompt": "Question text.",
          "choices": ["a", "b", "c"],
          "answer": "a",
          "explanation": "",
          "position": 1
        }
      ]
    }
  ],
  "vocabulary": [
    {
      "level_code": "VB",
      "lesson_number": 1,
      "category": "Days of the Week",
      "term": "Monday",
      "definition": "The first working day of the week.",
      "example": "I have class on Monday.",
      "position": 1
    }
  ]
}
```

## Rules

- `kind` is one of `mcq`, `fill_blank`, `true_false`.
- `mcq`: `answer` must exactly equal one of the `choices` strings.
- `fill_blank`: `choices` is `[]`; `answer` is the expected text. Grading normalizes
  case and whitespace. Multiple acceptable answers are written pipe-separated: `"books|book"`.
- `true_false`: `choices` is `["True", "False"]`; `answer` is `"True"` or `"False"`.
- `prompt`: mark a blank with three underscores `___`.
- `position`: 1-based order within the parent.
- `quizzes[].level_code` and `vocabulary[].level_code` reference `levels[].code`.
- `vocabulary[].lesson_number` is optional; when set it links the term to a lesson.
