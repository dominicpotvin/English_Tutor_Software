import { useParams } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import Breadcrumb from "../components/Breadcrumb";
import ExercisePlayer from "../components/ExercisePlayer";

export default function QuizPage() {
  const { id } = useParams();
  const state = useAsync(() => api.getQuiz(id!), [id]);

  return (
    <AsyncView state={state}>
      {({ quiz, exercises }) => (
        <div>
          <Breadcrumb
            items={[
              { label: "Dashboard", to: "/" },
              { label: "Quizzes", to: "/quizzes" },
              { label: quiz.title },
            ]}
          />
          <div className="page-head">
            <h1>{quiz.title}</h1>
            <p>{quiz.description}</p>
          </div>
          <ExercisePlayer exercises={exercises} />
        </div>
      )}
    </AsyncView>
  );
}
