import { useMemo } from "react";
import { Link } from "react-router-dom";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { getStoredUser } from "../lib/auth";

const doctorTips = [
  "Segarkan slot praktek secara berkala agar pasien tahu kapan bisa datang.",
  "Pantau permintaan booking dan ubah statusnya segera setelah diputuskan.",
  "Gunakan catatan singkat di kolom keluhan untuk mengingat konteks pasien.",
];

const patientTips = [
  "Gunakan fitur pencarian dokter untuk menemukan jadwal yang sesuai.",
  "Pilih tanggal yang masuk akal terhadap hari kerja dokter.",
  "Batalkan booking jika tidak jadi hadir agar slot bisa diberikan ke pasien lain.",
];

export default function DashboardPage() {
  const user = useMemo(() => getStoredUser(), []);
  const role = user?.role ?? "guest";
  const isDoctor = role === "doctor";
  const isPatient = role === "patient";

  const tips = isDoctor ? doctorTips : patientTips;

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <section className="bg-white border border-slate-200 rounded-lg p-4">
          <p className="text-sm text-slate-600 mb-1">Selamat datang</p>
          <h1 className="text-2xl font-semibold text-slate-900">
            {user?.name || "Pengguna"}{" "}
            {isDoctor && (
              <span className="text-sm font-normal text-slate-500">(Dokter)</span>
            )}
            {isPatient && (
              <span className="text-sm font-normal text-slate-500">(Pengguna)</span>
            )}
          </h1>
          <p className="text-sm text-slate-600 mt-2">
            {isDoctor
              ? "Jaga agar jadwal dan status booking selalu diperbarui supaya antrean tetap rapi."
              : isPatient
              ? "Booking dibuat langsung terhubung ke dokter, jadi cukup pilih jadwal yang masih tersedia."
              : "Masuk sebagai dokter atau pasien untuk melihat fitur secara penuh."}
          </p>
        </section>

        {(isDoctor || isPatient) && (
          <section className="bg-white border border-slate-200 rounded-lg p-4 space-y-3">
            <h2 className="text-lg font-semibold text-slate-900">Langkah berikutnya</h2>
            <div className="flex flex-wrap gap-3">
              {isDoctor && (
                <Link
                  to="/schedules"
                  className="px-4 py-2 rounded border border-slate-300 text-sm font-medium text-slate-700"
                >
                  Kelola Jadwal
                </Link>
              )}
              <Link
                to="/bookings"
                className="px-4 py-2 rounded border border-slate-300 text-sm font-medium text-slate-700"
              >
                Buka Halaman Booking
              </Link>
            </div>
            <ul className="list-disc pl-5 text-sm text-slate-600 space-y-1">
              {tips.map((tip) => (
                <li key={tip}>{tip}</li>
              ))}
            </ul>
          </section>
        )}

        {!isDoctor && !isPatient && (
          <section className="bg-white border border-slate-200 rounded-lg p-4">
            <p className="text-sm text-slate-600">
              Belum memiliki akun? Gunakan menu registrasi di halaman login dan pilih peran dokter atau pasien.
            </p>
          </section>
        )}
      </div>
    </DashboardLayout>
  );
}
