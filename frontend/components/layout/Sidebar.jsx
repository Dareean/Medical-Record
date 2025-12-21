import { Home, ClipboardList, LogOut } from "lucide-react";
import { useNavigate, useLocation } from "react-router-dom";
import { clearAuth, getStoredUser } from "../../src/lib/auth";
import Logo from "../layout/Logo";

export default function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  const role = getStoredUser()?.role;

  const menuItems = [
    { name: "Dashboard", icon: <Home size={18} />, path: "/dashboard" },
  ];

  if (role === "doctor") {
    menuItems.push({
      name: "Appointments",
      icon: <ClipboardList size={18} />,
      path: "/bookings",
    });
  }

  if (role === "patient") {
    menuItems.push({
      name: "My Appointments",
      icon: <ClipboardList size={18} />,
      path: "/bookings",
    });
  }

  return (
    <aside className="w-64 bg-white border-r border-slate-200 h-screen fixed left-0 top-0 hidden md:flex flex-col z-50">
      <div className="p-4 border-b border-slate-100">
        <Logo />
      </div>

      <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
        {menuItems.map((item) => {
          const isActive = location.pathname === item.path;
          return (
            <button
              key={item.name}
              onClick={() => navigate(item.path)}
              className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm ${
                isActive
                  ? "bg-red-50 text-red-600"
                  : "text-slate-600 hover:bg-slate-50"
              }`}
            >
              {item.icon}
              {item.name}
            </button>
          );
        })}
      </nav>

      <div className="p-4 border-t border-slate-100">
        <button
          onClick={() => {
            clearAuth();
            navigate("/login", { replace: true });
          }}
          className="w-full flex items-center gap-3 px-3 py-2 text-slate-600 hover:text-red-600 hover:bg-red-50 rounded-lg text-sm"
        >
          <LogOut size={20} />
          Sign Out
        </button>
      </div>
    </aside>
  );
}
