import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Mail, Lock, ArrowLeft, Loader2, Stethoscope } from "lucide-react";

export default function LoginPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  // GANTI URL INI DENGAN URL MOCKAPI ANDA SENDIRI
  // Contoh format: https://64a6...mockapi.io/api/v1/users
  const MOCK_API_URL = "https://64a6...mockapi.io/api/v1/users";

  const handleLogin = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    try {
      // 1. Simulasi delay network (opsional, biar kelihatan loadingnya)
      await new Promise((resolve) => setTimeout(resolve, 1500));

      // 2. Fetch data user dari Mock API
      // (Karena MockAPI tidak punya fitur Auth beneran, kita ambil semua user lalu filter manual)
      // const response = await fetch(MOCK_API_URL);
      // const users = await response.json();

      // LOGIKA DUMMY SEMENTARA (Supaya Anda bisa login tanpa set MockAPI dulu)
      // Hapus bagian ini jika sudah punya MockAPI URL
      const users = [
        { email: "admin@medicflow.id", password: "admin", name: "Dr. Admin" },
      ];

      // 3. Cek apakah email & password cocok
      const user = users.find(
        (u) => u.email === email && u.password === password
      );

      if (user) {
        // Login Berhasil!
        console.log("Login Success:", user);
        localStorage.setItem("user", JSON.stringify(user)); // Simpan sesi
        navigate("/dashboard"); // Redirect ke Dashboard (nanti kita buat)
      } else {
        // Login Gagal
        setError(
          "Invalid email or password. Try 'admin@medicflow.id' & 'admin'"
        );
      }
    } catch (err) {
      setError("Something went wrong. Please check your connection.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-white">
      {/* BAGIAN KIRI: FORM LOGIN */}
      <div className="w-full lg:w-1/2 p-8 md:p-12 flex flex-col justify-center relative">
        {/* Tombol Back */}
        <button
          onClick={() => navigate("/")}
          className="absolute top-8 left-8 flex items-center gap-2 text-slate-500 hover:text-red-500 transition-colors"
        >
          <ArrowLeft size={20} /> Back to Home
        </button>

        <div className="max-w-md w-full mx-auto space-y-8">
          <div className="text-center lg:text-left">
            <div className="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-red-50 text-red-500 mb-4">
              <Stethoscope size={24} />
            </div>
            <h2 className="text-3xl font-bold text-slate-900">Welcome Back</h2>
            <p className="text-slate-500 mt-2">
              Please enter your details to sign in.
            </p>
          </div>

          <form onSubmit={handleLogin} className="space-y-6">
            {/* Input Email */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-slate-700">
                Email
              </label>
              <div className="relative">
                <Mail
                  className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
                  size={20}
                />
                <input
                  type="email"
                  required
                  placeholder="Enter your email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full pl-10 pr-4 py-3 border border-slate-200 rounded-lg focus:ring-2 focus:ring-red-100 focus:border-red-400 outline-none transition-all"
                />
              </div>
            </div>

            {/* Input Password */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-slate-700">
                Password
              </label>
              <div className="relative">
                <Lock
                  className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
                  size={20}
                />
                <input
                  type="password"
                  required
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full pl-10 pr-4 py-3 border border-slate-200 rounded-lg focus:ring-2 focus:ring-red-100 focus:border-red-400 outline-none transition-all"
                />
              </div>
            </div>

            {/* Error Message */}
            {error && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
                {error}
              </div>
            )}

            {/* Submit Button */}
            <button
              type="submit"
              disabled={isLoading}
              className="w-full py-3 px-4 bg-red-500 text-white font-semibold rounded-lg shadow-lg hover:bg-red-600 focus:ring-4 focus:ring-red-200 transition-all disabled:opacity-70 flex justify-center items-center gap-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="animate-spin" size={20} /> Signing in...
                </>
              ) : (
                "Sign In"
              )}
            </button>
          </form>

          <p className="text-center text-sm text-slate-500">
            Don't have an account?{" "}
            <a href="#" className="text-red-500 font-semibold hover:underline">
              Contact Admin
            </a>
          </p>
        </div>
      </div>

      {/* BAGIAN KANAN: GAMBAR / ILUSTRASI */}
      <div className="hidden lg:block lg:w-1/2 bg-slate-50 relative overflow-hidden">
        <img
          src="https://images.unsplash.com/photo-1622253692010-333f2da6031d?q=80&w=1964&auto=format&fit=crop"
          alt="Doctor working"
          className="absolute inset-0 w-full h-full object-cover"
        />
        {/* Overlay Gradient */}
        <div className="absolute inset-0 bg-gradient-to-t from-slate-900/80 via-slate-900/20 to-transparent flex flex-col justify-end p-12 text-white">
          <h3 className="text-3xl font-bold mb-2">Modern Healthcare</h3>
          <p className="text-slate-200 text-lg">
            Manage your clinic with confidence and precision.
          </p>
        </div>
      </div>
    </div>
  );
}
