import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import Sidebar from "./Sidebar";
import { Search, User } from "lucide-react";
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

  const displayName = currentUser?.name || "Dr. Sarah";
  const displayRole = currentUser?.role
    ? currentUser.role.charAt(0).toUpperCase() + currentUser.role.slice(1)
    : "Cardiologist";

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar />

      <div className="flex-1 md:ml-64 transition-all">
        <header className="bg-white border-b border-slate-200 h-16 px-8 flex items-center justify-between sticky top-0 z-40">
          <div className="flex items-center gap-2 bg-slate-100 px-4 py-2 rounded-full w-96">
            <Search size={18} className="text-slate-400" />
            <input
              type="text"
              placeholder="Search patients, records..."
              className="bg-transparent border-none outline-none text-sm text-slate-700 w-full placeholder:text-slate-400"
            />
          </div>

          <div className="flex items-center gap-4">
            <div className="h-8 w-px bg-slate-300"></div>
            <div className="flex items-center gap-3">
              <div className="text-right hidden md:block">
                <p className="text-sm font-bold text-slate-700">
                  {displayName}
                </p>
                <p className="text-xs text-slate-500">{displayRole}</p>
              </div>
              <div className="w-10 h-10 bg-slate-200 rounded-full flex items-center justify-center text-slate-500 overflow-hidden border border-slate-300">
                <User size={20} />
              </div>
            </div>
          </div>
        </header>

        {/* Konten Halaman Berubah-ubah disini */}
        <main className="p-8">{children}</main>
      </div>
    </div>
  );
}
