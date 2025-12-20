import { useCallback, useEffect, useMemo, useState } from "react";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { doctorScheduleApi } from "../lib/api";
import { getStoredUser } from "../lib/auth";
import { CalendarClock, RefreshCw, Trash2, Edit3 } from "lucide-react";

const dayOptions = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

export default function AppointmentsPage() {
  const user = useMemo(() => getStoredUser(), []);
  const isDoctor = user?.role === "doctor";

  const [schedules, setSchedules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");
  const [editingId, setEditingId] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    work_day: "monday",
    start_time: "09:00",
    end_time: "11:00",
    patient_quota: 10,
  });

  const fetchSchedules = useCallback(async () => {
    if (!isDoctor) return;
    setLoading(true);
    setError("");
    try {
      const data = await doctorScheduleApi.getMine();
      setSchedules(data);
    } catch (err) {
      setError(err.message || "Failed to load schedules");
    } finally {
      setLoading(false);
    }
  }, [isDoctor]);

  useEffect(() => {
    if (!isDoctor) {
      setLoading(false);
      return;
    }
    fetchSchedules();
  }, [isDoctor, fetchSchedules]);

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const resetForm = () => {
    setFormData({
      work_day: "monday",
      start_time: "09:00",
      end_time: "11:00",
      patient_quota: 10,
    });
    setEditingId(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!isDoctor) return;
    setIsSubmitting(true);
    setError("");
    setStatus("");

    const payload = {
      work_day: formData.work_day,
      start_time: formData.start_time,
      end_time: formData.end_time,
      patient_quota: Number(formData.patient_quota) || 0,
    };

    try {
      if (editingId) {
        await doctorScheduleApi.update(editingId, payload);
        setStatus("Schedule updated");
      } else {
        await doctorScheduleApi.create(payload);
        setStatus("Schedule created");
      }
      resetForm();
      await fetchSchedules();
    } catch (err) {
      setError(err.message || "Failed to save schedule");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (schedule) => {
    setEditingId(schedule.id);
    setFormData({
      work_day: schedule.work_day,
      start_time: schedule.start_time?.slice(0, 5) || "09:00",
      end_time: schedule.end_time?.slice(0, 5) || "11:00",
      patient_quota: schedule.patient_quota ?? 0,
    });
  };

  const handleDelete = async (id) => {
    if (!window.confirm("Delete this schedule?")) return;
    setError("");
    setStatus("");
    try {
      await doctorScheduleApi.remove(id);
      setStatus("Schedule deleted");
      await fetchSchedules();
    } catch (err) {
      setError(err.message || "Failed to delete schedule");
    }
  };

  if (!isDoctor) {
    return (
      <DashboardLayout>
        <div className="bg-white border border-rose-100 rounded-2xl p-8 text-rose-600">
          This section is only available for doctor accounts.
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header className="flex flex-col gap-2">
          <p className="text-sm uppercase tracking-[0.3em] text-slate-500">
            Scheduling
          </p>
          <h1 className="text-3xl font-bold text-slate-900 flex items-center gap-3">
            <CalendarClock size={28} className="text-red-500" /> My Schedule
            Builder
          </h1>
          <p className="text-slate-500">
            Define availability windows and patient quota for each day.
          </p>
        </header>

        <section className="bg-white border border-slate-100 rounded-2xl p-6 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid md:grid-cols-4 gap-4">
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Work Day
                </label>
                <select
                  value={formData.work_day}
                  onChange={(e) => handleChange("work_day", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl"
                >
                  {dayOptions.map((day) => (
                    <option key={day} value={day}>
                      {day.charAt(0).toUpperCase() + day.slice(1)}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Start Time
                </label>
                <input
                  type="time"
                  value={formData.start_time}
                  onChange={(e) => handleChange("start_time", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl"
                />
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  End Time
                </label>
                <input
                  type="time"
                  value={formData.end_time}
                  onChange={(e) => handleChange("end_time", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl"
                />
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Patient Quota
                </label>
                <input
                  type="number"
                  min={0}
                  value={formData.patient_quota}
                  onChange={(e) =>
                    handleChange("patient_quota", e.target.value)
                  }
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl"
                />
              </div>
            </div>

            {error && (
              <div className="p-3 bg-red-50 text-red-600 rounded-xl border border-red-100 text-sm">
                {error}
              </div>
            )}
            {status && (
              <div className="p-3 bg-green-50 text-green-700 rounded-xl border border-green-100 text-sm">
                {status}
              </div>
            )}

            <div className="flex gap-4 flex-wrap">
              <button
                type="submit"
                disabled={isSubmitting}
                className="inline-flex items-center gap-2 px-6 py-3 bg-red-500 text-white rounded-xl font-semibold shadow hover:bg-red-600 disabled:opacity-60"
              >
                {editingId ? <Edit3 size={16} /> : <RefreshCw size={16} />}
                {editingId ? "Update Schedule" : "Add Schedule"}
              </button>
              {editingId && (
                <button
                  type="button"
                  onClick={resetForm}
                  className="inline-flex items-center gap-2 px-6 py-3 bg-slate-100 text-slate-700 rounded-xl font-semibold"
                >
                  Cancel
                </button>
              )}
            </div>
          </form>
        </section>

        <section className="bg-white border border-slate-100 rounded-2xl overflow-hidden">
          <div className="p-6 border-b border-slate-100 flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-slate-900">
                My Approved Slots
              </h2>
              <p className="text-sm text-slate-500">
                {schedules.length} schedule blocks
              </p>
            </div>
            <button
              onClick={fetchSchedules}
              className="inline-flex items-center gap-2 text-sm text-red-500 font-semibold"
            >
              <RefreshCw size={16} /> Refresh
            </button>
          </div>

          {loading ? (
            <p className="p-6 text-slate-500">Loading schedules...</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-semibold">
                  <tr>
                    <th className="px-6 py-4">Day</th>
                    <th className="px-6 py-4">Time</th>
                    <th className="px-6 py-4">Quota</th>
                    <th className="px-6 py-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {schedules.map((item) => (
                    <tr
                      key={item.id}
                      className="hover:bg-slate-50 transition-colors"
                    >
                      <td className="px-6 py-4 font-semibold text-slate-800">
                        {item.work_day.charAt(0).toUpperCase() +
                          item.work_day.slice(1)}
                      </td>
                      <td className="px-6 py-4">
                        {item.start_time?.slice(0, 5)} -{" "}
                        {item.end_time?.slice(0, 5)}
                      </td>
                      <td className="px-6 py-4">
                        {item.patient_quota || 0} patients
                      </td>
                      <td className="px-6 py-4 text-right space-x-2">
                        <button
                          className="inline-flex items-center gap-1 text-slate-500 hover:text-red-500"
                          onClick={() => handleEdit(item)}
                        >
                          <Edit3 size={16} /> Edit
                        </button>
                        <button
                          className="inline-flex items-center gap-1 text-rose-500"
                          onClick={() => handleDelete(item.id)}
                        >
                          <Trash2 size={16} /> Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {!schedules.length && (
                <p className="p-6 text-slate-500 text-center">
                  No schedules yet. Use the form above to create one.
                </p>
              )}
            </div>
          )}
        </section>
      </div>
    </DashboardLayout>
  );
}
