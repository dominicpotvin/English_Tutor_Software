import { Link, useParams } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import ProgressBar from "../components/ProgressBar";
import Breadcrumb from "../components/Breadcrumb";

export default function LevelPage() {
  const { id } = useParams();
  const state = useAsync(() => api.getLevel(id!), [id]);

  return (
    <AsyncView state={state}>
      {({ level, lessons, quizzes }) => (
        <div>
          <Breadcrumb items={[{ label: "Dashboard", to: "/" }, { label: level.name }]} />
          <div className="page-head">
            <h1>{level.name}</h1>
            <p>{level.description}</p>
          </div>

          <div className="section-title">Lessons</div>
          <div className="stack">
            {lessons.map((lesson) => (
              <Link key={lesson.id} to={`/lessons/${lesson.id}`} className="row-card">
                <div className="row-card-head">
                  <h3>
                    Lesson {lesson.number}: {lesson.title}
                  </h3>
                  {lesson.topicCount > 0 && (
                    <span className="badge">{lesson.topicCount} topics</span>
                  )}
                </div>
                <p className="summary">{lesson.summary}</p>
                {lesson.exerciseCount > 0 ? (
                  <ProgressBar value={lesson.masteredExercises} total={lesson.exerciseCount} />
                ) : (
                  <p className="vocab-hint">No practice content yet.</p>
                )}
              </Link>
            ))}
          </div>

          {quizzes.length > 0 && (
            <>
              <div className="section-title">Quizzes</div>
              <div className="stack">
                {quizzes.map((quiz) => (
                  <Link key={quiz.id} to={`/quizzes/${quiz.id}`} className="row-card">
                    <div className="row-card-head">
                      <h3>{quiz.title}</h3>
                      <span className="badge accent">{quiz.questionCount} questions</span>
                    </div>
                    <p className="summary">{quiz.description}</p>
                  </Link>
                ))}
              </div>
            </>
          )}
        </div>
      )}
    </AsyncView>
  );
}
