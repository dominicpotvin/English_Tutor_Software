import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import ProgressBar from "../components/ProgressBar";

function percent(part: number, whole: number): number {
  return whole > 0 ? Math.round((part / whole) * 100) : 0;
}

export default function ProgressPage() {
  const state = useAsync(() => api.getProgress(), []);

  return (
    <div>
      <div className="page-head">
        <h1>Your progress</h1>
        <p>A summary of your practice across the whole curriculum.</p>
      </div>
      <AsyncView state={state}>
        {({ levels, totals }) => (
          <div>
            <div className="stat-grid">
              <div className="stat">
                <div className="stat-value">
                  {totals.masteredExercises} / {totals.totalExercises}
                </div>
                <div className="stat-label">Exercises mastered</div>
              </div>
              <div className="stat">
                <div className="stat-value">{totals.totalAttempts}</div>
                <div className="stat-label">Answers submitted</div>
              </div>
              <div className="stat">
                <div className="stat-value">
                  {percent(totals.correctAttempts, totals.totalAttempts)}%
                </div>
                <div className="stat-label">Overall accuracy</div>
              </div>
            </div>

            <div className="section-title">By level</div>
            <div className="progress-table">
              {levels.map((lp) => (
                <div key={lp.levelId} className="progress-table-row">
                  <span className="name">{lp.name}</span>
                  <ProgressBar value={lp.masteredExercises} total={lp.totalExercises} />
                  <span className="badge">
                    {lp.totalAttempts > 0
                      ? `${percent(lp.correctAttempts, lp.totalAttempts)}% correct`
                      : "Not started"}
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}
      </AsyncView>
    </div>
  );
}
