import { useEffect, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { getStoredUser } from "../lib/auth";

export default function SchedulesPage() {
  const navigate = useNavigate();
  const user = useMemo(() => getStoredUser(), []);
  const isDoctor = user?.role === "doctor";

  useEffect(() => {
    if (isDoctor) {
      navigate("/bookings", { replace: true });
    }
  }, [isDoctor, navigate]);

  return (
    <DashboardLayout>
      <div className="space-y-4 bg-white border border-slate-200 rounded-xl p-6">
        <p className="text-xs uppercase tracking-[0.4em] text-red-400">
          Jadwal mandiri dimatikan
        </p>
        <h1 className="text-2xl font-semibold text-slate-900">
          Semua permintaan dibuat manual oleh pasien
        </h1>
        <p className="text-sm text-slate-600">
          Dokter tidak lagi perlu menetapkan kuota atau slot. Setiap permintaan
          datang langsung dari pasien dan kamu cukup menekan Approved atau
          Declined di halaman Booking Masuk.
        </p>
        <div className="flex flex-wrap gap-3">
          <button
            onClick={() => navigate("/bookings")}
            className="px-5 py-2 bg-red-500 text-white text-sm font-semibold rounded-lg"
          >
            Buka Booking Masuk
          </button>
          <button
            onClick={() => navigate("/dashboard")}
            className="px-5 py-2 border border-slate-200 text-sm font-semibold text-slate-600 rounded-lg"
          >
            Kembali ke Dashboard
          </button>
        </div>
      </div>
    </DashboardLayout>
  );
}
