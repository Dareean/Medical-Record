import { useEffect, useMemo, useState } from "react";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { doctorApi, doctorScheduleApi } from "../lib/api";
import { getStoredUser } from "../lib/auth";
import {
  Activity,
  AlertCircle,
  CalendarClock,
  Stethoscope,
  Users,
} from "lucide-react";

const dayOrder = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday",
];

const formatTime = (value) => (value ? value.slice(0, 5) : "-");
const capitalize = (value) =>
  value ? value.charAt(0).toUpperCase() + value.slice(1) : "-";
const specializationLabel = (doctor) =>
  doctor?.specialization?.name ||
  doctor?.specialization_id ||
  doctor?.["specialization_id "] ||
  "-";

export default function DashboardPage() {
  const user = useMemo(() => getStoredUser(), []);
  const role = user?.role || "";
  const isAdmin = role === "admin";
  const isDoctor = role === "doctor";

  const [doctors, setDoctors] = useState([]);
  const [schedules, setSchedules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      setError("");
      try {
        if (isAdmin) {
          const data = await doctorApi.getAll();
          setDoctors(data);
        } else if (isDoctor) {
          const data = await doctorScheduleApi.getMine();
          setSchedules(data);
        } else {
          setDoctors([]);
          setSchedules([]);
        }
      } catch (err) {
        setError(err.message || "Failed to load dashboard data");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [isAdmin, isDoctor]);

  const summaryCards = useMemo(() => {
    if (isAdmin) {
      const activeCount = doctors.filter(
        (doc) => doc.is_active !== false
      ).length;
      const specialtyCount = new Set(
        doctors
          .map((doc) => specializationLabel(doc))
          .filter((name) => name && name !== "-")
      ).size;

      return [
        {
          label: "Total Doctors",
          value: doctors.length,
          icon: Users,
          accent: "text-red-500 bg-red-50",
        },
        {
          label: "Active Licenses",
          value: activeCount,
          icon: Activity,
          accent: "text-emerald-500 bg-emerald-50",
        },
        {
          label: "Specialties",
          value: specialtyCount,
          icon: Stethoscope,
          accent: "text-slate-600 bg-slate-100",
        },
      ];
    }

    if (isDoctor) {
      const totalQuota = schedules.reduce(
        (sum, schedule) => sum + (schedule.patient_quota || 0),
        0
      );
      const uniqueDays = new Set(schedules.map((item) => item.work_day)).size;

      return [
        {
          label: "Weekly Sessions",
          value: schedules.length,
          icon: CalendarClock,
          accent: "text-red-500 bg-red-50",
        },
        {
          label: "Total Patient Quota",
          value: totalQuota,
          icon: Activity,
          accent: "text-emerald-500 bg-emerald-50",
        },
        {
          label: "Unique Days",
          value: uniqueDays,
          icon: Users,
          accent: "text-slate-600 bg-slate-100",
        },
      ];
    }

    return [];
  }, [isAdmin, isDoctor, doctors, schedules]);

  const topDoctors = useMemo(() => doctors.slice(0, 5), [doctors]);
  const sortedSchedules = useMemo(() => {
    if (!schedules?.length) return [];
    const order = dayOrder.reduce(
      (acc, day, idx) => ({ ...acc, [day]: idx }),
      {}
    );

    return [...schedules].sort((a, b) => {
      const dayDiff = (order[a.work_day] ?? 99) - (order[b.work_day] ?? 99);
      if (dayDiff !== 0) {
        return dayDiff;
      }
      return (a.start_time || "").localeCompare(b.start_time || "");
    });
  }, [schedules]);

  const heading = isAdmin
    ? "Clinic Command Center"
    : isDoctor
    ? "My Schedule Overview"
    : "Welcome";

  const description = isAdmin
    ? "Monitor doctor availability, licensing, and specialties in real time."
    : isDoctor
    ? "Track your approved slots and adjust quota before patients book."
    : "Sign in as admin or doctor to see live data.";

  return (
    <DashboardLayout>
      <div className="space-y-8">
        <header className="space-y-2">
          <p className="text-sm uppercase tracking-[0.4em] text-slate-500">
            Dashboard
          </p>
          <h1 className="text-3xl font-bold text-slate-900">{heading}</h1>
          <p className="text-slate-500">{description}</p>
        </header>

        {error && (
          <div className="flex items-center gap-3 rounded-2xl border border-rose-100 bg-rose-50 px-4 py-3 text-rose-600">
            <AlertCircle size={20} /> {error}
          </div>
        )}

        {summaryCards.length > 0 && (
          <section className="grid gap-4 md:grid-cols-3">
            {summaryCards.map((card) => (
              <article
                key={card.label}
                className="bg-white border border-slate-100 rounded-2xl p-5 flex items-center gap-4 shadow-sm"
              >
                <div
                  className={`w-12 h-12 rounded-2xl flex items-center justify-center ${card.accent}`}
                >
                  <card.icon size={20} />
                </div>
                <div>
                  <p className="text-sm text-slate-500">{card.label}</p>
                  <p className="text-2xl font-bold text-slate-900">
                    {card.value}
                  </p>
                </div>
              </article>
            ))}
          </section>
        )}

        <section className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
          <div className="p-6 border-b border-slate-100 flex justify-between items-center">
            <div>
              <h3 className="font-bold text-slate-900 text-lg">
                {isAdmin
                  ? "Doctor Snapshot"
                  : isDoctor
                  ? "Scheduled Sessions"
                  : "Overview"}
              </h3>
              <p className="text-sm text-slate-500">
                {isAdmin
                  ? `${doctors.length} registered doctors`
                  : isDoctor
                  ? `${schedules.length} slots approved`
                  : "Log in to access protected data"}
              </p>
            </div>
          </div>

          {loading ? (
            <p className="p-6 text-slate-500">Loading latest data...</p>
          ) : isAdmin ? (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-slate-500 uppercase font-semibold text-xs">
                  <tr>
                    <th className="px-6 py-4">Doctor</th>
                    <th className="px-6 py-4">Specialization</th>
                    <th className="px-6 py-4">Status</th>
                    <th className="px-6 py-4">License</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {topDoctors.map((doctor) => (
                    <tr
                      key={doctor.id}
                      className="hover:bg-slate-50 transition-colors"
                    >
                      <td className="px-6 py-4">
                        <p className="font-semibold text-slate-800">
                          {doctor?.user?.name || "-"}
                        </p>
                        <p className="text-xs text-slate-500">
                          {doctor?.user?.email}
                        </p>
                      </td>
                      <td className="px-6 py-4">
                        {specializationLabel(doctor)}
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-medium ${
                            doctor.is_active === false
                              ? "bg-slate-100 text-slate-600"
                              : "bg-emerald-100 text-emerald-700"
                          }`}
                        >
                          {doctor.is_active === false ? "Inactive" : "Active"}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        {doctor.license_number || "-"}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {!topDoctors.length && (
                <p className="p-6 text-slate-500 text-center">
                  No doctors registered yet.
                </p>
              )}
            </div>
          ) : isDoctor ? (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-slate-500 uppercase font-semibold text-xs">
                  <tr>
                    <th className="px-6 py-4">Day</th>
                    <th className="px-6 py-4">Time</th>
                    <th className="px-6 py-4">Quota</th>
                    <th className="px-6 py-4">Last Updated</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {sortedSchedules.map((schedule) => (
                    <tr
                      key={schedule.id}
                      className="hover:bg-slate-50 transition-colors"
                    >
                      <td className="px-6 py-4 font-semibold text-slate-800">
                        {capitalize(schedule.work_day)}
                      </td>
                      <td className="px-6 py-4">
                        {formatTime(schedule.start_time)} -{" "}
                        {formatTime(schedule.end_time)}
                      </td>
                      <td className="px-6 py-4">
                        {schedule.patient_quota || 0} patients
                      </td>
                      <td className="px-6 py-4 text-slate-500">
                        {new Date(
                          schedule.updated_at ||
                            schedule.created_at ||
                            Date.now()
                        ).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {!sortedSchedules.length && (
                <p className="p-6 text-slate-500 text-center">
                  No schedules configured. Use the Schedules page to create one.
                </p>
              )}
            </div>
          ) : (
            <p className="p-6 text-slate-500 text-center">
              Please log in with an admin or doctor account to see insights.
            </p>
          )}
        </section>
      </div>
    </DashboardLayout>
  );
}
