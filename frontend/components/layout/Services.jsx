import { Stethoscope, Activity, FileText, ArrowRight } from "lucide-react";

const servicesData = [
  {
    title: "General Checkup",
    desc: "Comprehensive health screening tailored to your age and needs.",
    icon: <Stethoscope size={24} />,
    color: "bg-blue-100 text-blue-600",
  },
  {
    title: "Lab Results",
    desc: "Fast and accurate laboratory results available directly on your dashboard.",
    icon: <Activity size={24} />,
    color: "bg-red-100 text-red-600",
  },
  {
    title: "Medical History",
    desc: "Secure storage for all your past diagnoses and prescriptions.",
    icon: <FileText size={24} />,
    color: "bg-green-100 text-green-600",
  },
];

export default function Services() {
  return (
    <section
      className="py-20 px-10 max-w-7xl mx-auto scroll-mt-32"
      id="services"
    >
      <div className="mb-12 text-center max-w-2xl mx-auto">
        <h2 className="text-red-500 font-bold tracking-wide uppercase text-sm mb-2">
          Our Services
        </h2>
        <h1 className="text-3xl md:text-4xl font-serif font-bold text-slate-900">
          High-quality services for your health
        </h1>
        <p className="text-slate-500 mt-4">
          We provide various features to make it easier for you to manage
          clinical data and patient health.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        {servicesData.map((service, index) => (
          <div
            key={index}
            className="group bg-white border border-gray-100 rounded-2xl p-8 shadow-sm hover:shadow-xl hover:-translate-y-1 transition-all duration-300"
          >
            <div
              className={`w-14 h-14 rounded-xl flex items-center justify-center mb-6 ${service.color}`}
            >
              {service.icon}
            </div>

            <h3 className="text-xl font-bold text-slate-800 mb-3">
              {service.title}
            </h3>
            <p className="text-slate-500 mb-6 leading-relaxed">
              {service.desc}
            </p>

            <a
              href="#"
              className="inline-flex items-center text-sm font-semibold text-slate-900 group-hover:text-red-400 transition-colors"
            >
              Learn More
              <ArrowRight
                size={16}
                className="ml-2 group-hover:translate-x-1 transition-transform"
              />
            </a>
          </div>
        ))}
      </div>
    </section>
  );
}
