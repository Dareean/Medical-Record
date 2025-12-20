import { useCallback, useEffect, useMemo, useState } from "react";
import DashboardLayout from "../../components/layout/DashboardLayout";
import { doctorApi } from "../lib/api";
import { getStoredUser } from "../lib/auth";
import { UserPlus, RefreshCw, Stethoscope, Trash2, Edit3 } from "lucide-react";

const initialForm = {
  name: "",
  email: "",
  specialization_id: "",
  gender: "male",
  address: "",
  license_number: "",
};

export default function DoctorsPage() {
  const user = useMemo(() => getStoredUser(), []);
  const isAdmin = user?.role === "admin";

  const [doctors, setDoctors] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [formData, setFormData] = useState(initialForm);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [status, setStatus] = useState("");
  const [editingId, setEditingId] = useState(null);

  const resolveSpecializationId = (doctor) =>
    doctor?.specialization_id ??
    doctor?.["specialization_id "] ??
    doctor?.specialization?.id ??
    "";

  const fetchDoctors = useCallback(async () => {
    if (!isAdmin) return;
    setLoading(true);
    setError("");
    try {
      const list = await doctorApi.getAll();
      setDoctors(list);
    } catch (err) {
      setError(err.message || "Failed to load doctors");
    } finally {
      setLoading(false);
    }
  }, [isAdmin]);

  useEffect(() => {
    fetchDoctors();
  }, [fetchDoctors]);

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const resetForm = () => {
    setFormData(initialForm);
    setEditingId(null);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setStatus("");
    setError("");

    const payload = {
      name: formData.name.trim(),
      email: formData.email.trim(),
      specialization_id: Number(formData.specialization_id),
      gender: formData.gender,
      address: formData.address.trim(),
      license_number: formData.license_number.trim(),
    };

    try {
      if (!payload.name || !payload.email || !payload.license_number) {
        throw new Error("Please complete the required fields");
      }
      if (Number.isNaN(payload.specialization_id)) {
        throw new Error("Specialization must be a number");
      }

      if (editingId) {
        await doctorApi.update(editingId, payload);
        setStatus("Doctor updated successfully");
      } else {
        await doctorApi.create(payload);
        setStatus("Doctor added successfully");
      }
      resetForm();
      await fetchDoctors();
    } catch (err) {
      setError(err.message || "Failed to save doctor");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleEdit = (doctor) => {
    setEditingId(doctor.id);
    setFormData({
      name: doctor?.user?.name || "",
      email: doctor?.user?.email || "",
      specialization_id: resolveSpecializationId(doctor)?.toString() || "",
      gender: doctor.gender || "male",
      address: doctor.address || "",
      license_number: doctor.license_number || "",
    });
  };

  const handleDelete = async (id) => {
    if (!window.confirm("Delete this doctor?")) return;
    setError("");
    setStatus("");
    try {
      await doctorApi.remove(id);
      setStatus("Doctor deleted");
      await fetchDoctors();
    } catch (err) {
      setError(err.message || "Failed to delete doctor");
    }
  };

  if (!isAdmin) {
    return (
      <DashboardLayout>
        <div className="bg-white border border-rose-100 text-rose-600 p-8 rounded-2xl">
          <p className="text-sm uppercase tracking-[0.3em] mb-2">Restricted</p>
          <h1 className="text-3xl font-bold text-slate-900 mb-3">Admin Only</h1>
          <p className="text-slate-600">
            Only administrators can manage doctors. Please log in with an admin
            account.
          </p>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header className="flex flex-col gap-2">
          <p className="text-sm uppercase tracking-[0.4em] text-slate-500">
            Doctors
          </p>
          <h1 className="text-3xl font-bold text-slate-900 flex items-center gap-3">
            <Stethoscope size={30} className="text-red-500" /> Doctor Management
          </h1>
          <p className="text-slate-500">
            Register new doctors, edit existing profiles, and keep licensing
            data up to date.
          </p>
        </header>

        <section className="bg-white border border-slate-100 rounded-2xl p-6 shadow-sm">
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Full Name
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => handleChange("name", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-red-100"
                  placeholder="Dr. Sarah Hart"
                  required
                />
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Email
                </label>
                <input
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleChange("email", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-red-100"
                  placeholder="doctor@hospital.com"
                  required
                />
              </div>
            </div>

            <div className="grid md:grid-cols-3 gap-4">
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Specialization ID
                </label>
                <input
                  type="number"
                  value={formData.specialization_id}
                  onChange={(e) =>
                    handleChange("specialization_id", e.target.value)
                  }
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-red-100"
                  placeholder="e.g. 1"
                  required
                />
                <p className="text-xs text-slate-400 mt-1">
                  Use IDs from the seeded specialization table.
                </p>
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  Gender
                </label>
                <select
                  value={formData.gender}
                  onChange={(e) => handleChange("gender", e.target.value)}
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl"
                >
                  <option value="male">Male</option>
                  <option value="female">Female</option>
                </select>
              </div>
              <div>
                <label className="text-sm font-semibold text-slate-600">
                  License Number
                </label>
                <input
                  type="text"
                  value={formData.license_number}
                  onChange={(e) =>
                    handleChange("license_number", e.target.value)
                  }
                  className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-red-100"
                  placeholder="STR-001"
                  required
                />
              </div>
            </div>

            <div>
              <label className="text-sm font-semibold text-slate-600">
                Address
              </label>
              <textarea
                value={formData.address}
                onChange={(e) => handleChange("address", e.target.value)}
                className="w-full mt-2 px-4 py-3 border border-slate-200 rounded-xl focus:ring-2 focus:ring-red-100"
                rows={3}
                placeholder="Clinic address or work location"
              />
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

            <div className="flex flex-wrap gap-4">
              <button
                type="submit"
                disabled={isSubmitting}
                className="inline-flex items-center gap-2 px-6 py-3 bg-red-500 text-white rounded-xl font-semibold shadow hover:bg-red-600 disabled:opacity-60"
              >
                <UserPlus size={18} />
                {editingId ? "Update Doctor" : "Add Doctor"}
              </button>
              {editingId && (
                <button
                  type="button"
                  onClick={resetForm}
                  className="inline-flex items-center gap-2 px-6 py-3 bg-slate-100 text-slate-700 rounded-xl font-semibold"
                >
                  <RefreshCw size={18} /> Cancel Edit
                </button>
              )}
            </div>
          </form>
        </section>

        <section className="bg-white rounded-2xl border border-slate-100">
          <div className="p-6 border-b border-slate-100 flex items-center justify-between">
            <div>
              <h2 className="text-xl font-semibold text-slate-900">
                Registered Doctors
              </h2>
              <p className="text-sm text-slate-500">
                {doctors.length} profiles
              </p>
            </div>
            <button
              onClick={fetchDoctors}
              className="text-sm text-red-500 font-semibold inline-flex items-center gap-2"
            >
              <RefreshCw size={16} /> Refresh
            </button>
          </div>

          {loading ? (
            <p className="p-6 text-slate-500">Loading doctors...</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-slate-600">
                <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-semibold">
                  <tr>
                    <th className="px-6 py-4">Doctor</th>
                    <th className="px-6 py-4">Specialization</th>
                    <th className="px-6 py-4">Gender</th>
                    <th className="px-6 py-4">STR</th>
                    <th className="px-6 py-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {doctors.map((doctor) => (
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
                        {doctor?.specialization?.name ||
                          resolveSpecializationId(doctor)}
                      </td>
                      <td className="px-6 py-4 capitalize">{doctor.gender}</td>
                      <td className="px-6 py-4">
                        {doctor.license_number || "-"}
                      </td>
                      <td className="px-6 py-4 text-right space-x-2">
                        <button
                          className="inline-flex items-center gap-1 text-slate-500 hover:text-red-500"
                          onClick={() => handleEdit(doctor)}
                        >
                          <Edit3 size={16} /> Edit
                        </button>
                        <button
                          className="inline-flex items-center gap-1 text-rose-500"
                          onClick={() => handleDelete(doctor.id)}
                        >
                          <Trash2 size={16} /> Delete
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {!doctors.length && (
                <p className="p-6 text-slate-500 text-center">
                  No doctors registered yet.
                </p>
              )}
            </div>
          )}
        </section>
      </div>
    </DashboardLayout>
  );
}
