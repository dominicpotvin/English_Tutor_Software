import { Routes, Route, Link } from "react-router-dom";
import Layout from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import LevelPage from "./pages/LevelPage";
import LessonPage from "./pages/LessonPage";
import PracticePage from "./pages/PracticePage";
import QuizListPage from "./pages/QuizListPage";
import QuizPage from "./pages/QuizPage";
import VocabularyPage from "./pages/VocabularyPage";
import ProgressPage from "./pages/ProgressPage";

function NotFound() {
  return (
    <div className="empty-state">
      <h2>Page not found</h2>
      <p>
        The page you are looking for does not exist. <Link to="/">Return to the dashboard</Link>.
      </p>
    </div>
  );
}

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/levels/:id" element={<LevelPage />} />
        <Route path="/lessons/:id" element={<LessonPage />} />
        <Route path="/topics/:id/practice" element={<PracticePage />} />
        <Route path="/quizzes" element={<QuizListPage />} />
        <Route path="/quizzes/:id" element={<QuizPage />} />
        <Route path="/vocabulary" element={<VocabularyPage />} />
        <Route path="/progress" element={<ProgressPage />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </Layout>
  );
}
