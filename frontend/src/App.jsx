import { createBrowserRouter, RouterProvider } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";
import SchedulesPage from "./pages/SchedulesPage";
import BookingsPage from "./pages/BookingsPage";
import LandingPage from "./pages/LandingPage";

const router = createBrowserRouter([
  {
    path: "/",
    element: <LoginPage />,
  },
  {
    path: "/landing",
    element: <LandingPage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/dashboard",
    element: <DashboardPage />,
  },
  {
    path: "/schedules",
    element: <SchedulesPage />,
  },
  {
    path: "/bookings",
    element: <BookingsPage />,
  },
  {
    path: "*",
    element: <NotFoundPage />,
  },
]);

function App() {
  return <RouterProvider router={router} />;
}

function NotFoundPage() {
  return (
    <main className="min-h-screen flex flex-col items-center justify-center bg-slate-50 px-6 text-center">
      <p className="text-xs uppercase tracking-[0.6em] text-red-400 mb-4">
        404
      </p>
      <h1 className="text-4xl md:text-5xl font-serif font-bold text-slate-900 mb-4">
        Halaman tidak ditemukan
      </h1>
      <p className="text-slate-600 max-w-2xl mb-8">
        URL ini belum dikonfigurasi. Gunakan tombol di bawah untuk kembali ke
        halaman yang terhubung dengan backend.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-4">
        <a
          href="/dashboard"
          className="px-6 py-3 rounded-full bg-red-500 text-white font-semibold hover:bg-red-600 transition-colors"
        >
          Buka Dashboard
        </a>
        <a
          href="/login"
          className="px-6 py-3 rounded-full border border-slate-200 text-slate-700 font-semibold hover:border-red-200 hover:text-red-500 transition-colors"
        >
          Kembali ke Login
        </a>
      </div>
    </main>
  );
}

export default App;
