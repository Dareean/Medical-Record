import { createBrowserRouter, RouterProvider } from "react-router-dom";
import LandingPage from "./pages/LandingPage";
import LoginPage from "./pages/LoginPage";
import DashboardPage from "./pages/DashboardPage";

const router = createBrowserRouter([
  {
    path: "/",
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
        URL yang Anda akses belum tersedia. Silakan kembali ke halaman utama
        atau buka menu login untuk melanjutkan.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-4">
        <a
          href="/"
          className="px-6 py-3 rounded-full bg-red-500 text-white font-semibold hover:bg-red-600 transition-colors"
        >
          Ke Landing Page
        </a>
        <a
          href="/login"
          className="px-6 py-3 rounded-full border border-slate-200 text-slate-700 font-semibold hover:border-red-200 hover:text-red-500 transition-colors"
        >
          Buka Login
        </a>
      </div>
    </main>
  );
}

export default App;
