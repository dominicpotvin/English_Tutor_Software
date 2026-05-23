import type {
  AttemptResult,
  LevelDetail,
  LevelSummary,
  LessonDetail,
  ProgressReport,
  QuizDetail,
  QuizSummary,
  TopicExercises,
  Vocabulary,
} from "./types";

const BASE = "/api";

async function http<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const body = (await res.json().catch(() => null)) as { error?: string } | null;
    throw new Error(body?.error ?? `Request failed (${res.status})`);
  }
  if (res.status === 204) {
    return undefined as T;
  }
  return (await res.json()) as T;
}

export const api = {
  listLevels: () => http<LevelSummary[]>("/levels"),
  getLevel: (id: string | number) => http<LevelDetail>(`/levels/${id}`),
  getLesson: (id: string | number) => http<LessonDetail>(`/lessons/${id}`),
  getTopicExercises: (id: string | number) => http<TopicExercises>(`/topics/${id}/exercises`),
  submitAttempt: (exerciseId: number, answer: string) =>
    http<AttemptResult>(`/exercises/${exerciseId}/attempt`, {
      method: "POST",
      body: JSON.stringify({ answer }),
    }),
  listQuizzes: () => http<QuizSummary[]>("/quizzes"),
  getQuiz: (id: string | number) => http<QuizDetail>(`/quizzes/${id}`),
  listVocabulary: () => http<Vocabulary[]>("/vocabulary"),
  getProgress: () => http<ProgressReport>("/progress"),
};
