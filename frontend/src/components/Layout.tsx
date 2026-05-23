import type { ReactNode } from "react";
import { Link, NavLink } from "react-router-dom";

export default function Layout({ children }: { children: ReactNode }) {
  return (
    <div className="app">
      <header className="app-header">
        <div className="app-header-inner">
          <Link to="/" className="brand">
            English Tutor
          </Link>
          <nav className="nav">
            <NavLink to="/" end>
              Dashboard
            </NavLink>
            <NavLink to="/quizzes">Quizzes</NavLink>
            <NavLink to="/vocabulary">Vocabulary</NavLink>
            <NavLink to="/progress">Progress</NavLink>
          </nav>
        </div>
      </header>
      <main className="app-main">{children}</main>
      <footer className="app-footer">
        English Tutor
      </footer>
    </div>
  );
}
