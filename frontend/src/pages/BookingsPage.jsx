import { useCallback, useEffect, useMemo, useState } from "react";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { doctorAppointmentApi, patientApi } from "../lib/api";
import { getStoredUser } from "../lib/auth";

const readableStatus = (status) => {
  switch (status) {
    case "Pending":
      return { label: "Menunggu", className: "text-amber-600" };
    case "Confirmed":
      return { label: "Disetujui", className: "text-emerald-600" };
    case "Rejected":
      return { label: "Ditolak", className: "text-slate-500" };
    case "Completed":
      return { label: "Selesai", className: "text-slate-600" };
    default:
      return { label: status || "-", className: "text-slate-600" };
  }
};

const toArray = (value) => (Array.isArray(value) ? value : []);

const DEFAULT_START_TIME = "09:00";

const getTodayISODate = () => {
  const now = new Date();
  const offsetMs = now.getTimezoneOffset() * 60000;
  return new Date(now.getTime() - offsetMs).toISOString().slice(0, 10);
};

const isValidTimeValue = (value) =>
  /^([01]\d|2[0-3]):[0-5]\d$/.test((value || "").slice(0, 5));

const buildLocalDateTime = (dateStr, timeStr) => {
  if (!dateStr || !timeStr) return null;
  const [year, month, day] = dateStr.split("-").map(Number);
  const [hour, minute] = timeStr.split(":").map(Number);
  if ([year, month, day, hour, minute].some((n) => Number.isNaN(n))) {
    return null;
  }
  return new Date(year, month - 1, day, hour, minute);
};

export default function BookingsPage() {
  const user = useMemo(() => getStoredUser(), []);
  const role = user?.role;

  if (role === "doctor") {
    return <DoctorBookings />;
  }

  return <PatientBookings />;
}

function PatientBookings() {
  const [searchTerm, setSearchTerm] = useState("");
  const [doctors, setDoctors] = useState([]);
  const [loadingDoctors, setLoadingDoctors] = useState(true);
  const [selectedDoctor, setSelectedDoctor] = useState(null);
  const todayIso = getTodayISODate();
  const [appointmentDate, setAppointmentDate] = useState(todayIso);
  const [startTime, setStartTime] = useState(DEFAULT_START_TIME);
  const [complaint, setComplaint] = useState("");
  const [history, setHistory] = useState([]);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const doctorOptions = toArray(doctors);
  const historyItems = toArray(history);
  const canSubmit = Boolean(selectedDoctor && appointmentDate && startTime);

  const handleAppointmentDateChange = (value) => {
    if (!value) {
      setAppointmentDate(getTodayISODate());
      return;
    }
    if (value < todayIso) {
      setAppointmentDate(todayIso);
      return;
    }
    setAppointmentDate(value);
  };

  const handleStartTimeChange = (value) => {
    const normalized = (value || "").slice(0, 5);
    if (!isValidTimeValue(normalized)) {
      setStartTime(DEFAULT_START_TIME);
      return;
    }
    setStartTime(normalized);
  };

  const fetchDoctors = useCallback(async () => {
    setLoadingDoctors(true);
    setError("");
    try {
      const list = await patientApi.searchDoctors(searchTerm);
      setDoctors(toArray(list));
    } catch (err) {
      setError(err.message || "Gagal memuat daftar dokter");
    } finally {
      setLoadingDoctors(false);
    }
  }, [searchTerm]);

  const fetchHistory = useCallback(async () => {
    try {
      const list = await patientApi.getAppointments();
      setHistory(toArray(list));
    } catch (err) {
      setError(err.message || "Gagal memuat riwayat booking");
    }
  }, []);

  useEffect(() => {
    fetchDoctors();
  }, [fetchDoctors]);

  useEffect(() => {
    fetchHistory();
  }, [fetchHistory]);

  const handleSelectDoctor = (doctor) => {
    setSelectedDoctor(doctor);
    setError("");
  };

  const handleBooking = async (event) => {
    event.preventDefault();
    setMessage("");

    if (!selectedDoctor || !appointmentDate || !startTime) {
      setError("Lengkapi tanggal dan jam yang diinginkan.");
      return;
    }

    const normalizedTime = startTime.slice(0, 5);
    if (!isValidTimeValue(normalizedTime)) {
      setError("Format jam tidak valid. Gunakan contoh 09:30.");
      return;
    }

    const desiredDateTime = buildLocalDateTime(appointmentDate, normalizedTime);
    if (!desiredDateTime) {
      setError("Tanggal atau jam tidak dapat diproses.");
      return;
    }

    if (desiredDateTime.getTime() <= Date.now()) {
      setError("Pilih tanggal dan jam yang berada di masa mendatang.");
      return;
    }

    setError("");
    setSubmitting(true);
    try {
      await patientApi.createAppointment({
        doctor_id: selectedDoctor.id,
        appointment_date: appointmentDate,
        start_time_slot: normalizedTime,
        complaint,
      });
      setMessage("Your Appointment has been sent to the Doctor");
      setComplaint("");
      setAppointmentDate(getTodayISODate());
      setStartTime(DEFAULT_START_TIME);
      await fetchHistory();
    } catch (err) {
      setError(err.message || "Failed to load booking");
    } finally {
      setSubmitting(false);
    }
  };

  const handleCancel = async (id) => {
    if (!window.confirm("Cancel your book?")) return;
    setError("");
    setMessage("");
    try {
      await patientApi.cancelAppointment(id);
      setMessage("Booking Canceled.");
      await fetchHistory();
    } catch (err) {
      setError(err.message || "Gagal membatalkan booking");
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <section className="bg-white border border-slate-200 rounded-lg p-4 space-y-3">
          <h1 className="text-xl font-semibold text-slate-900">
            Booking Dokter
          </h1>
          <p className="text-sm text-slate-600">
            Cari dokter, kirimkan tanggal dan jam yang kamu mau, lalu tunggu
            dokter menyetujui atau menolak permintaanmu.
          </p>

          <div className="flex flex-wrap gap-3">
            <input
              type="text"
              placeholder="Cari dokter atau spesialisasi"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="flex-1 min-w-[200px] px-3 py-2 border border-slate-300 rounded"
            />
            <button
              onClick={fetchDoctors}
              className="px-4 py-2 border border-slate-300 rounded text-sm text-slate-700"
            >
              Segarkan
            </button>
          </div>

          {loadingDoctors ? (
            <p className="text-sm text-slate-600">Memuat dokter...</p>
          ) : (
            <div className="space-y-2">
              {doctorOptions.length === 0 ? (
                <p className="text-sm text-slate-500">Belum ada data dokter.</p>
              ) : (
                doctorOptions.map((doctor) => (
                  <button
                    key={doctor.id}
                    onClick={() => handleSelectDoctor(doctor)}
                    className={`w-full text-left border rounded px-3 py-2 text-sm ${
                      selectedDoctor?.id === doctor.id
                        ? "border-red-400 bg-red-50"
                        : "border-slate-200"
                    }`}
                  >
                    <p className="font-semibold text-slate-900">
                      {doctor?.user?.name || "Tanpa nama"}
                    </p>
                    <p className="text-xs text-slate-500">
                      {doctor?.specialization?.name || `Dokter #${doctor.id}`}
                    </p>
                  </button>
                ))
              )}
            </div>
          )}
        </section>

        {selectedDoctor && (
          <section className="bg-white border border-slate-200 rounded-lg p-4 space-y-4">
            <div className="flex flex-col gap-6 lg:flex-row">
              <div className="flex-1 space-y-4">
                <div>
                  <p className="text-sm uppercase tracking-wide text-slate-500">
                    {selectedDoctor?.specialization?.name || "Dokter"}
                  </p>
                  <h2 className="text-xl font-semibold text-slate-900">
                    Ajukan jadwal sesuai kebutuhanmu
                  </h2>
                  <p className="text-sm text-slate-600">
                    Tentukan tanggal dan jam yang diinginkan. Dokter akan
                    meninjau lalu menyetujui atau menolak sesuai ketersediaan
                    mereka.
                  </p>
                </div>

                <form onSubmit={handleBooking} className="space-y-3">
                  <label className="text-sm text-slate-600 flex flex-col gap-1">
                    Tanggal kunjungan
                    <input
                      type="date"
                      value={appointmentDate}
                      min={todayIso}
                      onChange={(e) =>
                        handleAppointmentDateChange(e.target.value)
                      }
                      className="w-full px-3 py-2 border border-slate-300 rounded text-sm"
                    />
                  </label>
                  <label className="text-sm text-slate-600 flex flex-col gap-1">
                    Jam preferensi
                    <input
                      type="time"
                      value={startTime}
                      onChange={(e) => handleStartTimeChange(e.target.value)}
                      className="w-full px-3 py-2 border border-slate-300 rounded text-sm"
                    />
                  </label>
                  <label className="text-sm text-slate-600 flex flex-col gap-1">
                    Keluhan atau catatan (opsional)
                    <textarea
                      value={complaint}
                      onChange={(e) => setComplaint(e.target.value)}
                      className="w-full px-3 py-2 border border-slate-300 rounded text-sm"
                      rows={3}
                      placeholder="Contoh: kontrol rutin, konsultasi nutrisi"
                    />
                  </label>
                  <button
                    type="submit"
                    className="px-4 py-2 bg-red-500 text-white rounded text-sm font-medium disabled:opacity-50"
                    disabled={!canSubmit || submitting}
                  >
                    {submitting ? "Mengirim..." : "Kirim permintaan"}
                  </button>
                </form>
              </div>

              <div className="flex-1 space-y-3 bg-slate-50 border border-slate-200 rounded-lg p-4">
                <h3 className="text-sm font-semibold text-slate-800">
                  Bagaimana prosesnya?
                </h3>
                <p className="text-sm text-slate-600">
                  Kamu bebas mengusulkan jadwal. Dokter akan meninjaunya dan
                  hanya perlu memilih approved atau declined tanpa mengatur
                  kuota apa pun.
                </p>
                <ul className="space-y-2 text-sm text-slate-600">
                  <li className="flex gap-2">
                    <span className="text-red-500">1.</span>
                    Kirim tanggal dan jam yang kamu inginkan lengkap dengan
                    keluhan singkat.
                  </li>
                  <li className="flex gap-2">
                    <span className="text-red-500">2.</span>
                    Dokter menerima permintaanmu di panel booking, lalu memilih
                    Approved / Declined.
                  </li>
                  <li className="flex gap-2">
                    <span className="text-red-500">3.</span>
                    Kamu akan mendapatkan status terbaru di daftar riwayat
                    booking.
                  </li>
                </ul>
              </div>
            </div>
          </section>
        )}

        <section className="bg-white border border-slate-200 rounded-lg p-4 space-y-3">
          <h2 className="text-lg font-semibold text-slate-900">
            Riwayat Booking
          </h2>
          {historyItems.length === 0 ? (
            <p className="text-sm text-slate-600">Belum ada data booking.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-left">
                <thead className="text-slate-500">
                  <tr>
                    <th className="py-2">Dokter</th>
                    <th className="py-2">Tanggal</th>
                    <th className="py-2">Jam</th>
                    <th className="py-2">Status</th>
                    <th className="py-2 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="text-slate-700">
                  {historyItems.map((item) => {
                    const badge = readableStatus(item.status);
                    const canCancel =
                      item.status === "Pending" || item.status === "Confirmed";
                    return (
                      <tr key={item.id} className="border-t border-slate-100">
                        <td className="py-2">
                          {item?.doctor?.user?.name ||
                            `Dokter #${item.doctor_id}`}
                        </td>
                        <td className="py-2">
                          {item.appointment_date?.slice(0, 10) || "-"}
                        </td>
                        <td className="py-2">
                          {item.start_time_slot?.slice(0, 5) || "-"}
                        </td>
                        <td className={`py-2 ${badge.className}`}>
                          {badge.label}
                        </td>
                        <td className="py-2 text-right">
                          {canCancel && (
                            <button
                              onClick={() => handleCancel(item.id)}
                              className="text-sm text-red-500"
                            >
                              Batalkan
                            </button>
                          )}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </section>

        {error && (
          <div className="bg-rose-50 border border-rose-200 text-rose-700 px-3 py-2 rounded text-sm">
            {error}
          </div>
        )}
        {message && (
          <div className="bg-emerald-50 border border-emerald-200 text-emerald-700 px-3 py-2 rounded text-sm">
            {message}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}

function DoctorBookings() {
  const [requests, setRequests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [updatingId, setUpdatingId] = useState(null);
  const requestItems = toArray(requests);

  const fetchRequests = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const list = await doctorAppointmentApi.getMine();
      setRequests(toArray(list));
    } catch (err) {
      setError(err.message || "Gagal memuat booking");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchRequests();
  }, [fetchRequests]);

  const handleUpdate = async (id, status) => {
    setUpdatingId(id);
    setError("");
    setMessage("");
    try {
      await doctorAppointmentApi.updateStatus(id, status);
      setMessage("Status booking diperbarui.");
      await fetchRequests();
    } catch (err) {
      setError(err.message || "Gagal memperbarui status");
    } finally {
      setUpdatingId(null);
    }
  };

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <section className="bg-white border border-slate-200 rounded-lg p-4 space-y-2">
          <h1 className="text-xl font-semibold text-slate-900">
            Permintaan Booking
          </h1>
          <p className="text-sm text-slate-600">
            Terima atau tolak booking dari pasien langsung dari daftar berikut.
          </p>
        </section>

        <section className="bg-white border border-slate-200 rounded-lg p-4">
          {loading ? (
            <p className="text-sm text-slate-600">Memuat data...</p>
          ) : requestItems.length === 0 ? (
            <p className="text-sm text-slate-600">Belum ada booking.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-left">
                <thead className="text-slate-500">
                  <tr>
                    <th className="py-2">Pasien</th>
                    <th className="py-2">Tanggal</th>
                    <th className="py-2">Jam</th>
                    <th className="py-2">Keluhan</th>
                    <th className="py-2">Status</th>
                    <th className="py-2 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="text-slate-700">
                  {requestItems.map((item) => {
                    const badge = readableStatus(item.status);
                    const disabled =
                      item.status === "Rejected" || item.status === "Completed";
                    return (
                      <tr key={item.id} className="border-t border-slate-100">
                        <td className="py-2">
                          {item?.patient?.name || `Pasien #${item.patient_id}`}
                        </td>
                        <td className="py-2">
                          {item.appointment_date?.slice(0, 10) || "-"}
                        </td>
                        <td className="py-2">
                          {item.start_time_slot?.slice(0, 5) || "-"}
                        </td>
                        <td className="py-2 text-slate-500">
                          {item.complaint || "-"}
                        </td>
                        <td className={`py-2 ${badge.className}`}>
                          {badge.label}
                        </td>
                        <td className="py-2 text-right space-x-2">
                          <button
                            onClick={() => handleUpdate(item.id, "Confirmed")}
                            disabled={disabled || updatingId === item.id}
                            className="text-sm text-emerald-600 disabled:opacity-40"
                          >
                            Setujui
                          </button>
                          <button
                            onClick={() => handleUpdate(item.id, "Rejected")}
                            disabled={disabled || updatingId === item.id}
                            className="text-sm text-rose-600 disabled:opacity-40"
                          >
                            Tolak
                          </button>
                          <button
                            onClick={() => handleUpdate(item.id, "Completed")}
                            disabled={disabled || updatingId === item.id}
                            className="text-sm text-slate-600 disabled:opacity-40"
                          >
                            Selesai
                          </button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </section>

        {error && (
          <div className="bg-rose-50 border border-rose-200 text-rose-700 px-3 py-2 rounded text-sm">
            {error}
          </div>
        )}
        {message && (
          <div className="bg-emerald-50 border border-emerald-200 text-emerald-700 px-3 py-2 rounded text-sm">
            {message}
          </div>
        )}
      </div>
    </DashboardLayout>
  );
}
