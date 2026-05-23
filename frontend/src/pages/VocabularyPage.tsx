import { useState } from "react";
import { api } from "../lib/api";
import { useAsync } from "../lib/useAsync";
import AsyncView from "../components/AsyncView";
import type { Vocabulary } from "../lib/types";

function categoriesOf(items: Vocabulary[]): string[] {
  const seen: string[] = [];
  for (const v of items) {
    if (v.category && !seen.includes(v.category)) {
      seen.push(v.category);
    }
  }
  return seen;
}

export default function VocabularyPage() {
  const state = useAsync(() => api.listVocabulary(), []);
  const [category, setCategory] = useState("All");
  const [study, setStudy] = useState(false);
  const [revealed, setRevealed] = useState<Set<number>>(new Set());

  function toggleReveal(id: number) {
    setRevealed((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  return (
    <div>
      <div className="page-head">
        <h1>Vocabulary</h1>
        <p>Browse key words and expressions. Turn on study mode to test your memory.</p>
      </div>
      <AsyncView state={state}>
        {(items) => {
          const categories = ["All", ...categoriesOf(items)];
          const visible =
            category === "All" ? items : items.filter((v) => v.category === category);
          return (
            <div>
              <div className="toolbar">
                <div className="filter-row">
                  {categories.map((c) => (
                    <button
                      key={c}
                      className={`chip${c === category ? " active" : ""}`}
                      onClick={() => setCategory(c)}
                    >
                      {c}
                    </button>
                  ))}
                </div>
                <button
                  className={`chip${study ? " active" : ""}`}
                  onClick={() => setStudy(!study)}
                >
                  Study mode
                </button>
              </div>

              {visible.length === 0 ? (
                <div className="empty-state">
                  <p>No vocabulary in this category yet.</p>
                </div>
              ) : (
                <div className="vocab-grid">
                  {visible.map((v) => {
                    const open = !study || revealed.has(v.id);
                    return (
                      <div
                        key={v.id}
                        className={`vocab-card${study ? " clickable" : ""}`}
                        onClick={() => {
                          if (study) {
                            toggleReveal(v.id);
                          }
                        }}
                      >
                        <div className="vocab-term">{v.term}</div>
                        {open ? (
                          <>
                            <div className="vocab-def">{v.definition}</div>
                            {v.example && <div className="vocab-example">{v.example}</div>}
                          </>
                        ) : (
                          <div className="vocab-hint">Click to reveal</div>
                        )}
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          );
        }}
      </AsyncView>
    </div>
  );
}
