import { useParams } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import Breadcrumb from "../components/Breadcrumb";
import ExercisePlayer from "../components/ExercisePlayer";

export default function PracticePage() {
  const { id } = useParams();
  const state = useAsync(() => api.getTopicExercises(id!), [id]);

  return (
    <AsyncView state={state}>
      {({ topic, lesson, exercises }) => (
        <div>
          <Breadcrumb
            items={[
              { label: "Dashboard", to: "/" },
              { label: `Lesson ${lesson.number}`, to: `/lessons/${lesson.id}` },
              { label: topic.title },
            ]}
          />
          <div className="page-head">
            <h1>{topic.title}</h1>
            <p>Answer each exercise and check your work as you go.</p>
          </div>
          <ExercisePlayer exercises={exercises} />
        </div>
      )}
    </AsyncView>
  );
}
