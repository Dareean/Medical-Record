import { useNavigate } from "react-router-dom";

export default function DashboardPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col items-center justify-center px-6 text-center">
      <p className="text-sm uppercase tracking-[0.3em] text-slate-400 mb-4">
        Coming Soon
      </p>
      <h1 className="text-4xl md:text-5xl font-serif font-bold text-slate-900 mb-6">
        Dashboard MedicFlow
      </h1>
      <p className="text-slate-600 max-w-2xl mb-8">
        Kami sedang mempersiapkan halaman dashboard dengan rangkaian fitur
        analitik pasien, ringkasan catatan medis, dan penjadwalan cerdas. Untuk
        sementara Anda bisa kembali melihat landing page atau mengakses menu
        lainnya.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-4">
        <button
          onClick={() => navigate("/")}
          className="px-6 py-3 rounded-full bg-red-500 text-white font-semibold hover:bg-red-600 transition-colors"
        >
          Kembali ke Landing
        </button>
        <button
          onClick={() => navigate("/login")}
          className="px-6 py-3 rounded-full border border-slate-200 text-slate-700 font-semibold hover:border-red-200 hover:text-red-500 transition-colors"
        >
          Masuk ke Login
        </button>
      </div>
    </div>
  );
}
