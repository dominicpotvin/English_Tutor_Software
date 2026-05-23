interface Props {
  value: number;
  total: number;
}

/** ProgressBar shows a value/total ratio as a bar with a numeric label. */
export default function ProgressBar({ value, total }: Props) {
  const pct = total > 0 ? Math.round((value / total) * 100) : 0;
  const complete = total > 0 && value >= total;
  return (
    <div className="progress-line">
      <div className="progress">
        <div
          className={`progress-fill${complete ? " is-complete" : ""}`}
          style={{ width: `${pct}%` }}
        />
      </div>
      <span className="count">
        {value} / {total}
      </span>
    </div>
  );
}
