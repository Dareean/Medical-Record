import DashboardLayout from "../../components/layout/DashboardLayout";
import { UserPlus, Search, ChevronRight } from "lucide-react";

export default function PatientsPage() {
  const patients = [
    {
      name: "Siti Aminah",
      mrn: "MF-1023",
      lastVisit: "04 Dec 2025",
      status: "Active",
    },
    {
      name: "Budi Santoso",
      mrn: "MF-0998",
      lastVisit: "02 Dec 2025",
      status: "Scheduled",
    },
    {
      name: "Rina Wati",
      mrn: "MF-0875",
      lastVisit: "28 Nov 2025",
      status: "Follow-up",
    },
    {
      name: "Joko Anwar",
      mrn: "MF-0811",
      lastVisit: "15 Nov 2025",
      status: "Archived",
    },
  ];

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <p className="text-sm text-slate-500 uppercase tracking-[0.3em]">
              Patients
            </p>
            <h1 className="text-3xl font-bold text-slate-900">
              Patient Registry
            </h1>
            <p className="text-slate-500">
              Manage, search, and onboard new patients.
            </p>
          </div>
          <button className="inline-flex items-center gap-2 bg-red-500 text-white px-6 py-3 rounded-xl shadow hover:bg-red-600 transition-colors">
            <UserPlus size={18} /> Add Patient
          </button>
        </header>

        <div className="bg-white border border-slate-100 rounded-2xl p-4 flex items-center gap-3">
          <Search size={18} className="text-slate-400" />
          <input
            type="text"
            placeholder="Search by name, MRN, or phone"
            className="flex-1 bg-transparent outline-none text-sm text-slate-600"
          />
        </div>

        <section className="bg-white rounded-2xl border border-slate-100">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-semibold">
              <tr>
                <th className="px-6 py-4">Patient</th>
                <th className="px-6 py-4">MRN</th>
                <th className="px-6 py-4">Last Visit</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {patients.map((patient) => (
                <tr
                  key={patient.mrn}
                  className="hover:bg-slate-50 transition-colors"
                >
                  <td className="px-6 py-4 font-semibold text-slate-800">
                    {patient.name}
                  </td>
                  <td className="px-6 py-4">{patient.mrn}</td>
                  <td className="px-6 py-4">{patient.lastVisit}</td>
                  <td className="px-6 py-4">
                    <span className="px-3 py-1 rounded-full text-xs font-medium bg-slate-100 text-slate-600">
                      {patient.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <button className="inline-flex items-center gap-1 text-red-500 font-medium">
                      View <ChevronRight size={16} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </section>
      </div>
    </DashboardLayout>
  );
}
