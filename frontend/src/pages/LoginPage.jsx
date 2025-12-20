import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Mail, Lock, ArrowLeft, Loader2, Stethoscope } from "lucide-react";
import { authApi } from "../lib/api";
import { saveAuth } from "../lib/auth";

export default function LoginPage() {
  const navigate = useNavigate();
  const [authMode, setAuthMode] = useState("login");
  const isRegisterMode = authMode === "register";
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState("patient");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");
    setStatus("");

    try {
      if (isRegisterMode) {
        await authApi.register({ name, email, password, role });
        setStatus("Account created. You can now sign in.");
        setAuthMode("login");
        setPassword("");
        setName("");
        setRole("patient");
      } else {
        const result = await authApi.login({ email, password });
        if (result?.data) {
          saveAuth(result.data);
        }
        navigate("/dashboard");
      }
    } catch (err) {
      setError(err.message || "Something went wrong");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex bg-white">
      <div className="w-full p-8 md:p-12 flex flex-col justify-center relative">
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
            <h2 className="text-3xl font-bold text-slate-900">
              {isRegisterMode ? "Create Account" : "Welcome Back"}
            </h2>
            <p className="text-slate-500 mt-2">
              {isRegisterMode
                ? "Register new medical staff access."
                : "Please enter your details to sign in."}
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {isRegisterMode && (
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Full Name
                </label>
                <input
                  type="text"
                  required
                  placeholder="Enter full name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full px-4 py-3 border border-slate-200 rounded-lg focus:ring-2 focus:ring-red-100 focus:border-red-400 outline-none transition-all"
                />
              </div>
            )}
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

            {isRegisterMode && (
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700">
                  Role
                </label>
                <select
                  value={role}
                  onChange={(e) => setRole(e.target.value)}
                  className="w-full px-4 py-3 border border-slate-200 rounded-lg focus:ring-2 focus:ring-red-100 focus:border-red-400 outline-none transition-all text-sm"
                >
                  <option value="doctor">Doctor</option>
                  <option value="patient">Patient</option>
                </select>
              </div>
            )}

            {error && (
              <div className="p-3 bg-red-50 text-red-600 text-sm rounded-lg border border-red-100">
                {error}
              </div>
            )}

            {status && (
              <div className="p-3 bg-green-50 text-green-700 text-sm rounded-lg border border-green-100">
                {status}
              </div>
            )}

            <button
              type="submit"
              disabled={isLoading}
              className="w-full py-3 px-4 bg-red-500 text-white font-semibold rounded-lg shadow-lg hover:bg-red-600 focus:ring-4 focus:ring-red-200 transition-all disabled:opacity-70 flex justify-center items-center gap-2"
            >
              {isLoading ? (
                <>
                  <Loader2 className="animate-spin" size={20} />
                  {isRegisterMode ? "Submitting..." : "Signing in..."}
                </>
              ) : isRegisterMode ? (
                "Register"
              ) : (
                "Sign In"
              )}
            </button>
          </form>

          <p className="text-center text-sm text-slate-500">
            {isRegisterMode
              ? "Already have an account?"
              : "Don't have an account?"}{" "}
            <button
              type="button"
              onClick={() => {
                setAuthMode(isRegisterMode ? "login" : "register");
                setError("");
                setStatus("");
                if (!isRegisterMode) {
                  setName("");
                  setRole("patient");
                }
              }}
              className="text-red-500 font-semibold hover:underline"
            >
              {isRegisterMode ? "Sign in" : "Register"}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
}
