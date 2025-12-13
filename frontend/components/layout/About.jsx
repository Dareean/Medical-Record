import { ShieldCheck, Heart, Clock, Users } from "lucide-react";

export default function About() {
  const stats = [
    { label: "Partner Clinics", value: "50+" },
    { label: "Patients Served", value: "10k+" },
    { label: "Doctors Online", value: "19M+" },
    { label: "Satisfaction", value: "97%" },
  ];

  const values = [
    {
      icon: <ShieldCheck className="w-6 h-6 text-red-600" />,
      title: "Data Security",
      desc: "We use the highest encryption standards to maintain the confidentiality of your patient records.",
    },
    {
      icon: <Clock className="w-6 h-6 text-red-600" />,
      title: "Real-time Efficiency",
      desc: "No more manual queues. All schedules and data are updated instantly (real-time).",
    },
    {
      icon: <Heart className="w-6 h-6 text-red-600" />,
      title: "Patient Centric",
      desc: "A specially designed interface that allows medical personnel to focus more on patient care.",
    },
  ];

  return (
    <section id="about" className="py-20 px-25 bg-slate-50 scroll-mt-32">
      <div className="max-w-7xl mx-auto px-6 md:px-12">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-12 items-center mb-20">
          <div className="relative">
            <img
              src="https://images.unsplash.com/photo-1576091160399-112ba8d25d1d?ixlib=rb-1.2.1&auto=format&fit=crop&w=1050&q=80"
              alt="Medical Team"
              className="relative rounded-2xl shadow-2xl w-full object-cover h-[450px]"
            />
          </div>
          <div className="space-y-6">
            <h2 className="text-red-500 font-bold tracking-wide uppercase text-sm">
              About Us
            </h2>
            <h1 className="text-4xl md:text-5xl font-serif font-bold text-slate-900 leading-tight">
              Simplifying healthcare, <br />
              <span className="text-red-400">one click at a time.</span>
            </h1>
            <p className="text-slate-600 text-lg ">
              We built this platform with one goal in mind: to eliminate the
              complexity of hospital administration. Doctors should be busy
              saving lives, not searching for missing files.
            </p>
            <p className="text-slate-600 text-lg">
              Our system combines high-level security with a user-friendly
              interface, making clinic management an enjoyable experience.
            </p>
          </div>
        </div>

        <div className="bg-white rounded-2xl p-8 shadow-sm border border-slate-100 grid grid-cols-2 md:grid-cols-4 gap-8 text-center mb-20">
          {stats.map((stat, index) => (
            <div key={index} className="space-y-1">
              <p className="text-3xl font-bold text-slate-800">{stat.value}</p>
              <p className="text-sm text-slate-500 font-medium uppercase tracking-wider">
                {stat.label}
              </p>
            </div>
          ))}
        </div>

        <div className="text-center mb-12">
          <h2 className="text-2xl font-bold text-slate-900">Why Choose Us?</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {values.map((item, index) => (
            <div
              key={index}
              className="bg-white p-8 rounded-xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="w-12 h-12 bg-red-50 rounded-lg flex items-center justify-center mb-4">
                {item.icon}
              </div>
              <h3 className="text-xl font-bold text-slate-800 mb-2">
                {item.title}
              </h3>
              <p className="text-slate-500 leading-relaxed">{item.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
