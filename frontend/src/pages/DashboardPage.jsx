import DashboardLayout from "../../components/layout/DashboardLayout";
import {
  Users,
  UserPlus,
  Calendar,
  DollarSign,
  TrendingUp,
  MoreHorizontal,
} from "lucide-react";

export default function DashboardPage() {
  const todayPatients = [
    {
      name: "Siti Aminah",
      time: "09:00 AM",
      type: "Check-up",
      status: "Completed",
    },
    {
      name: "Budi Santoso",
      time: "10:30 AM",
      type: "Consultation",
      status: "In Progress",
    },
    {
      name: "Rina Wati",
      time: "11:45 AM",
      type: "Lab Test",
      status: "Waiting",
    },
    {
      name: "Joko Anwar",
      time: "01:00 PM",
      type: "Surgery",
      status: "Scheduled",
    },
  ];

  return (
    <DashboardLayout>
      <div className="space-y-8">
        {/* Header Section */}
        <div>
          <h1 className="text-2xl font-bold text-slate-800">
            Dashboard Overview
          </h1>
          <p className="text-slate-500">
            Welcome back, Dr. Sarah. Here's what's happening today.
          </p>
        </div>

        {/* Recent Activity / Today's Schedule */}
        <div className="bg-white rounded-2xl shadow-sm border border-slate-100 overflow-hidden">
          <div className="p-6 border-b border-slate-100 flex justify-between items-center">
            <h3 className="font-bold text-slate-800 text-lg">
              Today's Appointments
            </h3>
            <button className="text-sm text-red-500 font-medium hover:underline">
              View All
            </button>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm text-slate-600">
              <thead className="bg-slate-50 text-slate-500 uppercase font-semibold text-xs">
                <tr>
                  <th className="px-6 py-4">Patient Name</th>
                  <th className="px-6 py-4">Time</th>
                  <th className="px-6 py-4">Type</th>
                  <th className="px-6 py-4">Status</th>
                  <th className="px-6 py-4 text-right">Action</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {todayPatients.map((patient, i) => (
                  <tr key={i} className="hover:bg-slate-50 transition-colors">
                    <td className="px-6 py-4 font-medium text-slate-800">
                      {patient.name}
                    </td>
                    <td className="px-6 py-4">{patient.time}</td>
                    <td className="px-6 py-4">{patient.type}</td>
                    <td className="px-6 py-4">
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium 
                        ${
                          patient.status === "Completed"
                            ? "bg-green-100 text-green-700"
                            : patient.status === "In Progress"
                            ? "bg-blue-100 text-blue-700"
                            : patient.status === "Waiting"
                            ? "bg-orange-100 text-orange-700"
                            : "bg-slate-100 text-slate-600"
                        }`}
                      >
                        {patient.status}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-right">
                      <button className="text-slate-400 hover:text-slate-600">
                        <MoreHorizontal size={20} />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
}
