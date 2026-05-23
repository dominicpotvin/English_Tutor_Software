import { Link } from "react-router-dom";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import ProgressBar from "../components/ProgressBar";

export default function Dashboard() {
  const state = useAsync(() => api.listLevels(), []);

  return (
    <div>
      <div className="page-head">
        <h1>Your English curriculum</h1>
        <p>Work through each level at your own pace. Select a level to begin.</p>
      </div>
      <AsyncView state={state}>
        {(levels) => (
          <div className="card-grid">
            {levels.map((level) => (
              <Link key={level.id} to={`/levels/${level.id}`} className="tile">
                <span className="tile-eyebrow">
                  {level.code} &middot; {level.lessonCount} lessons
                </span>
                <h3>{level.name}</h3>
                <p>{level.description}</p>
                {level.totalExercises > 0 ? (
                  <ProgressBar value={level.masteredExercises} total={level.totalExercises} />
                ) : (
                  <p className="vocab-hint">Curriculum outline ready.</p>
                )}
              </Link>
            ))}
          </div>
        )}
      </AsyncView>
    </div>
  );
}
