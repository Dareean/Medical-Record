import DashboardLayout from "../../components/layout/DashboardLayout";
import { FileText, Download } from "lucide-react";

export default function RecordsPage() {
  const records = [
    {
      id: "REC-2031",
      patient: "Siti Aminah",
      type: "Cardiology",
      updated: "04 Dec 2025",
      size: "1.2 MB",
    },
    {
      id: "REC-2032",
      patient: "Budi Santoso",
      type: "Radiology",
      updated: "30 Nov 2025",
      size: "3.8 MB",
    },
    {
      id: "REC-2033",
      patient: "Rina Wati",
      type: "Laboratory",
      updated: "29 Nov 2025",
      size: "980 KB",
    },
    {
      id: "REC-2034",
      patient: "Joko Anwar",
      type: "Surgery",
      updated: "22 Nov 2025",
      size: "5.4 MB",
    },
  ];

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header>
          <p className="text-sm uppercase tracking-[0.3em] text-slate-500">
            Records
          </p>
          <h1 className="text-3xl font-bold text-slate-900">
            Medical Records Library
          </h1>
          <p className="text-slate-500">
            Access signed reports, imaging results, and lab summaries.
          </p>
        </header>

        <section className="bg-white rounded-2xl border border-slate-100">
          <table className="w-full text-left text-sm text-slate-600">
            <thead className="bg-slate-50 text-slate-500 uppercase text-xs font-semibold">
              <tr>
                <th className="px-6 py-4">Record ID</th>
                <th className="px-6 py-4">Patient</th>
                <th className="px-6 py-4">Discipline</th>
                <th className="px-6 py-4">Last Updated</th>
                <th className="px-6 py-4">Size</th>
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
                      {record.id}
                    </span>
                  </td>
                  <td className="px-6 py-4">{record.patient}</td>
                  <td className="px-6 py-4">{record.type}</td>
                  <td className="px-6 py-4">{record.updated}</td>
                  <td className="px-6 py-4">{record.size}</td>
                  <td className="px-6 py-4 text-right">
                    <button className="inline-flex items-center gap-2 text-red-500 font-medium">
                      <Download size={16} /> Download
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
