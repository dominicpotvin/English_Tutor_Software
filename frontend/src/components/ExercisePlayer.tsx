import { useState } from "react";
import type { Exercise } from "../lib/types";
import { api } from "../lib/api";

/** Renders an exercise prompt, drawing a blank wherever "___" appears. */
function PromptText({ text }: { text: string }) {
  const parts = text.split("___");
  return (
    <p className="exercise-prompt">
      {parts.map((part, i) => (
        <span key={i}>
          {part}
          {i < parts.length - 1 && <span className="blank" />}
        </span>
      ))}
    </p>
  );
}

interface AnswerState {
  correct: boolean;
  correctAnswer: string;
  explanation: string;
}

/**
 * ExercisePlayer walks the learner through a set of exercises one at a time,
 * grading each answer through the API and showing a final score.
 */
export default function ExercisePlayer({ exercises }: { exercises: Exercise[] }) {
  const [index, setIndex] = useState(0);
  const [choice, setChoice] = useState("");
  const [text, setText] = useState("");
  const [answer, setAnswer] = useState<AnswerState | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [outcomes, setOutcomes] = useState<boolean[]>([]);

  function reset() {
    setChoice("");
    setText("");
    setAnswer(null);
    setError("");
  }

  function restart() {
    setIndex(0);
    setOutcomes([]);
    reset();
  }

  async function check() {
    const exercise = exercises[index];
    const given = exercise.kind === "fill_blank" ? text.trim() : choice;
    if (!given || submitting) {
      return;
    }
    setSubmitting(true);
    setError("");
    try {
      const result = await api.submitAttempt(exercise.id, given);
      setAnswer({
        correct: result.correct,
        correctAnswer: result.correctAnswer,
        explanation: result.explanation,
      });
      setOutcomes((prev) => [...prev, result.correct]);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Could not submit your answer.");
    } finally {
      setSubmitting(false);
    }
  }

  function next() {
    setIndex((i) => i + 1);
    reset();
  }

  if (exercises.length === 0) {
    return (
      <div className="empty-state">
        <p>There are no exercises here yet.</p>
      </div>
    );
  }

  if (index >= exercises.length) {
    const score = outcomes.filter(Boolean).length;
    return (
      <div className="exercise-summary card">
        <h2>Set complete</h2>
        <p className="exercise-score">
          {score} / {exercises.length}
        </p>
        <button className="btn btn-primary" onClick={restart}>
          Practice again
        </button>
      </div>
    );
  }

  const exercise = exercises[index];
  const hasChoices = exercise.kind === "mcq" || exercise.kind === "true_false";
  const given = exercise.kind === "fill_blank" ? text.trim() : choice;
  const answered = answer !== null;
  const isLast = index + 1 >= exercises.length;

  return (
    <div className="exercise-card">
      <div className="exercise-progress">
        <span>
          Exercise {index + 1} of {exercises.length}
        </span>
        <div className="progress">
          <div
            className="progress-fill"
            style={{ width: `${(index / exercises.length) * 100}%` }}
          />
        </div>
      </div>

      <PromptText text={exercise.prompt} />

      {hasChoices && (
        <div className="choices">
          {exercise.choices.map((option) => {
            let className = "choice";
            if (answered) {
              if (option === answer.correctAnswer) {
                className += " correct";
              } else if (option === choice) {
                className += " wrong";
              }
            } else if (option === choice) {
              className += " selected";
            }
            return (
              <button
                key={option}
                className={className}
                disabled={answered}
                onClick={() => setChoice(option)}
              >
                {option}
              </button>
            );
          })}
        </div>
      )}

      {exercise.kind === "fill_blank" && (
        <input
          className="fill-input"
          value={text}
          disabled={answered}
          placeholder="Type your answer"
          onChange={(e) => setText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              void check();
            }
          }}
        />
      )}

      {answered && (
        <div className={`feedback ${answer.correct ? "correct" : "wrong"}`}>
          <strong>{answer.correct ? "Correct" : "Not quite"}</strong>
          {!answer.correct && (
            <p>
              Correct answer: <span className="answer">{answer.correctAnswer}</span>
            </p>
          )}
          {answer.explanation && <p className="feedback-explanation">{answer.explanation}</p>}
        </div>
      )}

      {error && <p className="status status-error">{error}</p>}

      <div className="exercise-actions">
        {answered ? (
          <button className="btn btn-primary" onClick={next}>
            {isLast ? "See results" : "Next exercise"}
          </button>
        ) : (
          <button
            className="btn btn-primary"
            disabled={!given || submitting}
            onClick={() => void check()}
          >
            {submitting ? "Checking..." : "Check answer"}
          </button>
        )}
      </div>
    </div>
  );
}
