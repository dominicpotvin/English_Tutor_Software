import type { ReactNode } from "react";

interface AsyncLike<T> {
  data?: T;
  error?: string;
  loading: boolean;
}

interface Props<T> {
  state: AsyncLike<T>;
  children: (data: T) => ReactNode;
}

/** AsyncView renders loading and error states, then hands data to children. */
export default function AsyncView<T>({ state, children }: Props<T>) {
  if (state.loading) {
    return <div className="status">Loading...</div>;
  }
  if (state.error) {
    return <div className="status status-error">Could not load this content: {state.error}</div>;
  }
  if (state.data === undefined) {
    return <div className="status status-error">No content available.</div>;
  }
  return <>{children(state.data)}</>;
}
