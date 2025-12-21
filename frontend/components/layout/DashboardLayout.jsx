import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "./Sidebar";
import { getStoredUser } from "../../src/lib/auth";

export default function DashboardLayout({ children }) {
  const navigate = useNavigate();
  const [currentUser, setCurrentUser] = useState(() => getStoredUser());

  useEffect(() => {
    if (!currentUser) {
      navigate("/login", { replace: true });
    }
  }, [currentUser, navigate]);

  useEffect(() => {
    if (typeof window === "undefined") {
      return undefined;
    }

    const handleStorage = () => {
      setCurrentUser(getStoredUser());
    };

    window.addEventListener("storage", handleStorage);
    return () => window.removeEventListener("storage", handleStorage);
  }, []);

  const displayName = currentUser?.name || "Pengguna";
  const displayRole = currentUser?.role
    ? currentUser.role.charAt(0).toUpperCase() + currentUser.role.slice(1)
    : "Tanpa Peran";

  if (!currentUser) {
    return null;
  }

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar />

      <div className="flex-1 md:ml-64 transition-all">
        <header className="bg-white border-b border-slate-200 px-6 py-4 sticky top-0 z-40">
          <p className="text-xs uppercase tracking-[0.3em] text-slate-400">
            You Are Logged As
          </p>
          <div className="mt-2">
            <p className="text-lg font-semibold text-slate-900">
              {displayName}
            </p>
            <p className="text-sm text-slate-500">{displayRole}</p>
          </div>
        </header>

        <main className="p-6">{children}</main>
      </div>
    </div>
  );
}
