import {
  Home,
  Users,
  Calendar,
  LogOut,
  FileText,
  Activity,
} from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { clearAuth } from "../../src/lib/auth";

export default function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    { name: "Dashboard", icon: <Home size={20} />, path: "/dashboard" },
    { name: "Doctors", icon: <Users size={20} />, path: "/doctors" },
    {
      name: "Schedules",
      icon: <Calendar size={20} />,
      path: "/appointments",
    },
    { name: "Records", icon: <FileText size={20} />, path: "/records" },
  ];

  return (
    <aside className="w-64 bg-white border-r border-slate-200 h-screen fixed left-0 top-0 hidden md:flex flex-col z-50">
      {/* Logo Area */}
      <div className="p-6 border-b border-slate-100 flex items-center gap-3">
        <div className="w-8 h-8 bg-red-500 rounded-lg flex items-center justify-center text-white">
          <Activity size={20} />
        </div>
        <span className="text-xl font-bold text-slate-800">MedicFlow</span>
      </div>

      {/* Menu Items */}
      <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
        {menuItems.map((item) => {
          const isActive = location.pathname === item.path;
          return (
            <button
              key={item.name}
              onClick={() => navigate(item.path)}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-medium text-sm ${
                isActive
                  ? "bg-red-50 text-red-600 shadow-sm"
                  : "text-slate-500 hover:bg-slate-50 hover:text-slate-900"
              }`}
            >
              {item.icon}
              {item.name}
            </button>
          );
        })}
      </nav>

      {/* Logout Area */}
      <div className="p-4 border-t border-slate-100">
        <button
          onClick={() => {
            clearAuth();
            localStorage.removeItem("user");
            navigate("/login", { replace: true });
          }}
          className="w-full flex items-center gap-3 px-4 py-3 text-slate-500 hover:text-red-600 hover:bg-red-50 rounded-xl transition-all text-sm font-medium"
        >
          <LogOut size={20} />
          Sign Out
        </button>
      </div>
    </aside>
  );
}
