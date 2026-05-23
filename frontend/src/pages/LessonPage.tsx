import { Link, useParams } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import ProgressBar from "../components/ProgressBar";
import Breadcrumb from "../components/Breadcrumb";
import Markdown from "../components/Markdown";

export default function LessonPage() {
  const { id } = useParams();
  const state = useAsync(() => api.getLesson(id!), [id]);

  return (
    <AsyncView state={state}>
      {({ lesson, level, topics }) => (
        <div>
          <Breadcrumb
            items={[
              { label: "Dashboard", to: "/" },
              { label: level.name, to: `/levels/${level.id}` },
              { label: `Lesson ${lesson.number}` },
            ]}
          />
          <div className="page-head">
            <h1>
              Lesson {lesson.number}: {lesson.title}
            </h1>
            <p>{lesson.summary}</p>
          </div>

          {topics.length === 0 ? (
            <div className="empty-state">
              <h2>No content yet</h2>
              <p>
                This lesson does not have practice content yet. It can be added through the
                MCP connector.
              </p>
            </div>
          ) : (
            topics.map((topic) => (
              <div key={topic.id} className="topic">
                <div className="topic-head">
                  <h2>{topic.title}</h2>
                  <span className="badge">{topic.exerciseCount} exercises</span>
                </div>
                <Markdown content={topic.explanation} />
                <div className="topic-foot">
                  <ProgressBar value={topic.masteredExercises} total={topic.exerciseCount} />
                  <Link className="btn btn-primary" to={`/topics/${topic.id}/practice`}>
                    Practice
                  </Link>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </AsyncView>
  );
}
