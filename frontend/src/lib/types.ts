// TypeScript mirror of the backend JSON contract.

export interface Level {
  id: number;
  code: string;
  name: string;
  description: string;
  position: number;
}

export interface LevelSummary extends Level {
  lessonCount: number;
  totalExercises: number;
  masteredExercises: number;
}

export interface Lesson {
  id: number;
  levelId: number;
  number: number;
  title: string;
  summary: string;
  position: number;
}

export interface LessonSummary extends Lesson {
  topicCount: number;
  exerciseCount: number;
  masteredExercises: number;
}

export interface Topic {
  id: number;
  lessonId: number;
  title: string;
  explanation: string;
  position: number;
}

export interface TopicSummary extends Topic {
  exerciseCount: number;
  masteredExercises: number;
}

export interface Quiz {
  id: number;
  levelId: number | null;
  title: string;
  description: string;
  position: number;
}

export interface QuizSummary extends Quiz {
  questionCount: number;
  masteredCount: number;
}

export type ExerciseKind = "mcq" | "fill_blank" | "true_false";

export interface Exercise {
  id: number;
  topicId: number | null;
  quizId: number | null;
  kind: ExerciseKind;
  prompt: string;
  choices: string[];
  answer: string;
  explanation: string;
  position: number;
}

export interface Vocabulary {
  id: number;
  levelId: number | null;
  lessonId: number | null;
  category: string;
  term: string;
  definition: string;
  example: string;
  position: number;
}

export interface AttemptResult {
  attemptId: number;
  correct: boolean;
  correctAnswer: string;
  explanation: string;
}

export interface LevelProgress {
  levelId: number;
  code: string;
  name: string;
  totalExercises: number;
  masteredExercises: number;
  totalAttempts: number;
  correctAttempts: number;
}

export interface ProgressTotals {
  totalExercises: number;
  masteredExercises: number;
  totalAttempts: number;
  correctAttempts: number;
}

export interface ProgressReport {
  levels: LevelProgress[];
  totals: ProgressTotals;
}

// Composite responses returned by detail endpoints.

export interface LevelDetail {
  level: Level;
  lessons: LessonSummary[];
  quizzes: QuizSummary[];
}

export interface LessonDetail {
  lesson: Lesson;
  level: Level;
  topics: TopicSummary[];
}

export interface TopicExercises {
  topic: Topic;
  lesson: Lesson;
  exercises: Exercise[];
}

export interface QuizDetail {
  quiz: Quiz;
  exercises: Exercise[];
}
