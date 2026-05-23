import { useCallback, useEffect, useState } from "react";

interface AsyncState<T> {
  data?: T;
  error?: string;
  loading: boolean;
}

/**
 * useAsync runs an async fetcher and tracks its state. The fetcher re-runs
 * whenever `deps` change; `reload` re-runs it on demand.
 */
export function useAsync<T>(fetcher: () => Promise<T>, deps: unknown[]) {
  const [state, setState] = useState<AsyncState<T>>({ loading: true });

  const run = useCallback(() => {
    let active = true;
    setState({ loading: true });
    fetcher()
      .then((data) => {
        if (active) setState({ data, loading: false });
      })
      .catch((err: unknown) => {
        if (active) {
          setState({
            error: err instanceof Error ? err.message : "Unexpected error",
            loading: false,
          });
        }
      });
    return () => {
      active = false;
    };
    // The fetcher is intentionally re-created by callers from `deps`.
  }, deps);

  useEffect(() => run(), [run]);

  return { ...state, reload: run };
}
