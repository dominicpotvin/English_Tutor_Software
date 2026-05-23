import { Link } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import ProgressBar from "../components/ProgressBar";

export default function QuizListPage() {
  const state = useAsync(() => api.listQuizzes(), []);

  return (
    <div>
      <div className="page-head">
        <h1>Quizzes</h1>
        <p>Test what you have learned with a full assessment.</p>
      </div>
      <AsyncView state={state}>
        {(quizzes) =>
          quizzes.length === 0 ? (
            <div className="empty-state">
              <p>No quizzes are available yet.</p>
            </div>
          ) : (
            <div className="card-grid">
              {quizzes.map((quiz) => (
                <Link key={quiz.id} to={`/quizzes/${quiz.id}`} className="tile">
                  <span className="tile-eyebrow">{quiz.questionCount} questions</span>
                  <h3>{quiz.title}</h3>
                  <p>{quiz.description}</p>
                  <ProgressBar value={quiz.masteredCount} total={quiz.questionCount} />
                </Link>
              ))}
            </div>
          )
        }
      </AsyncView>
    </div>
  );
}
