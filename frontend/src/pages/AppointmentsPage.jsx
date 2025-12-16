import DashboardLayout from "../../components/layout/DashboardLayout";
import { Calendar, Clock, CheckCircle2 } from "lucide-react";

export default function AppointmentsPage() {
  const schedule = [
    {
      patient: "Siti Aminah",
      time: "09:00",
      type: "General Check-up",
      room: "A-01",
      status: "Checked-in",
    },
    {
      patient: "Budi Santoso",
      time: "10:30",
      type: "Cardio Consult",
      room: "B-02",
      status: "On Site",
    },
    {
      patient: "Rina Wati",
      time: "13:15",
      type: "Lab Review",
      room: "Lab",
      status: "Confirmed",
    },
    {
      patient: "Joko Anwar",
      time: "15:00",
      type: "Post-Op",
      room: "C-04",
      status: "Pending",
    },
  ];

  return (
    <DashboardLayout>
      <div className="space-y-6">
        <header className="flex flex-col gap-2">
          <p className="text-sm uppercase tracking-[0.3em] text-slate-500">
            Today
          </p>
          <h1 className="text-3xl font-bold text-slate-900 flex items-center gap-3">
            <Calendar size={28} className="text-red-500" /> Appointment Center
          </h1>
          <p className="text-slate-500">
            Monitor upcoming visits, assign rooms, and track patient arrivals.
          </p>
        </header>

        <section className="grid gap-6 md:grid-cols-2">
          {schedule.map((item) => (
            <article
              key={item.patient}
              className="bg-white border border-slate-100 rounded-2xl p-6 flex flex-col gap-4 shadow-sm"
            >
              <div className="flex items-start justify-between">
                <div>
                  <h3 className="text-xl font-semibold text-slate-900">
                    {item.patient}
                  </h3>
                  <p className="text-sm text-slate-500">{item.type}</p>
                </div>
                <span className="px-3 py-1 rounded-full bg-slate-100 text-xs font-semibold text-slate-600">
                  {item.status}
                </span>
              </div>

              <div className="flex items-center gap-4 text-sm text-slate-600">
                <span className="inline-flex items-center gap-2">
                  <Clock size={16} /> {item.time}
                </span>
                <span className="inline-flex items-center gap-2">
                  <CheckCircle2 size={16} /> Room {item.room}
                </span>
              </div>
            </article>
          ))}
        </section>
      </div>
    </DashboardLayout>
  );
}
