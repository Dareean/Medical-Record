import { useEffect, useMemo, useState } from "react";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { FileText, Download, AlertCircle } from "lucide-react";
import { doctorApi, doctorScheduleApi } from "../lib/api";
import { getStoredUser } from "../lib/auth";

const formatDate = (value) => {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleDateString();
};

const specializationName = (doctor) =>
  doctor?.specialization?.name ||
  doctor?.specialization_id ||
  doctor?.["specialization_id "] ||
  "-";

export default function RecordsPage() {
  const user = useMemo(() => getStoredUser(), []);
  const role = user?.role ?? "";
  const isAdmin = role === "admin";
  const isDoctor = role === "doctor";

  const [records, setRecords] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const load = async () => {
      setLoading(true);
      setError("");
      try {
        if (isAdmin) {
          const data = await doctorApi.getAll();
          setRecords(data);
        } else if (isDoctor) {
          const data = await doctorScheduleApi.getMine();
          setRecords(data);
        } else {
          setRecords([]);
        }
      } catch (err) {
        setError(err.message || "Failed to fetch records");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [isAdmin, isDoctor]);

  const handleDownload = (record, label) => {
    const blob = new Blob([JSON.stringify(record, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${label}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  const heading = isAdmin
    ? "Doctor Records"
    : isDoctor
    ? "Schedule Archive"
    : "Records";
  const subtext = isAdmin
    ? "Download signed doctor profiles and licensing data."
    : isDoctor
    ? "Keep a copy of your approved availability."
    : "Sign in to access protected records.";

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header>
          <p className="text-sm uppercase tracking-[0.3em] text-slate-500">
            Records
          </p>
          <h1 className="text-3xl font-bold text-slate-900">{heading}</h1>
          <p className="text-slate-500">{subtext}</p>
        </header>

        {error && (
          <div className="flex items-center gap-3 rounded-2xl border border-rose-100 bg-rose-50 px-4 py-3 text-rose-600">
            <AlertCircle size={20} /> {error}
          </div>
        )}

        <section className="bg-white rounded-2xl border border-slate-100">
          {loading ? (
            <p className="p-6 text-slate-500">Loading records...</p>
          ) : (isAdmin || isDoctor) && records.length ? (
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-semibold">
                <tr>
                  <th className="px-6 py-4">Reference</th>
                  <th className="px-6 py-4">
                    {isAdmin ? "Specialization" : "Time"}
                  </th>
                  <th className="px-6 py-4">{isAdmin ? "License" : "Quota"}</th>
                  <th className="px-6 py-4">Last Updated</th>
                  <th className="px-6 py-4 text-right">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {records.map((record) => (
                  <tr
                    key={record.id}
                    className="hover:bg-slate-50 transition-colors"
                  >
                    <td className="px-6 py-4 font-semibold text-slate-800">
                      <span className="inline-flex items-center gap-2">
                        <FileText size={16} className="text-red-500" />
                        {isAdmin
                          ? record?.user?.name || record.id
                          : record.work_day}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      {isAdmin
                        ? specializationName(record)
                        : `${record.start_time?.slice(
                            0,
                            5
                          )} - ${record.end_time?.slice(0, 5)}`}
                    </td>
                    <td className="px-6 py-4">
                      {isAdmin
                        ? record.license_number || "-"
                        : `${record.patient_quota || 0} patients`}
                    </td>
                    <td className="px-6 py-4">
                      {formatDate(record.updated_at || record.created_at)}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button
                        className="inline-flex items-center gap-2 text-red-500 font-medium"
                        onClick={() =>
                          handleDownload(
                            record,
                            isAdmin
                              ? record?.user?.name || `doctor-${record.id}`
                              : `${record.work_day}-${record.id}`
                          )
                        }
                      >
                        <Download size={16} /> Download JSON
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <p className="p-6 text-slate-500 text-center">
              {isAdmin || isDoctor
                ? "No records available yet."
                : "Please log in with the correct role to view records."}
            </p>
          )}
        </section>
      </div>
    </DashboardLayout>
  );
}
