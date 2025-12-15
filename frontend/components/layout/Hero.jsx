import { ArrowRight } from "lucide-react";

export default function Hero() {
  return (
    <main className="max-w-7xl mx-auto px-6 lg:px-10 my-30">
      <section
        id="hero"
        className="flex flex-col lg:flex-row items-center justify-between gap-10 lg:gap-16"
      >
        <div className="space-y-6 max-w-xl">
          <div className="border border-red-500 w-fit px-3 py-1 rounded-full text-red-500 font-medium">
            <p className="">#1 Favorite</p>
          </div>
          <h1 className="text-7xl font-black font-serif">
            Streamline your clinical{" "}
            <span className="text-red-400">workflow.</span>
          </h1>
          <h2 className="text-gray-500 text-lg">
            Effortlessly handle patient records, history, and appointments in
            one place. Less paperwork, more time for patients.
          </h2>
          <div className="flex items-center gap-4 text-sm">
            <button className="group flex items-center gap-2 text-white bg-red-400 px-10 py-3 rounded-full cursor-pointer hover:bg-red-300 color-transition">
              Get Started
              <ArrowRight
                size={17}
                className="group-hover:translate-x-1 transition-transform"
              />
            </button>
          </div>
        </div>

        <div className="p-5 space-y-2 w-full max-w-2xl">
          <img
            className="w-full h-auto object-cover"
            src="/dashboard-medical-record.png"
            alt="Medical Illustration"
          />
        </div>
      </section>
    </main>
  );
}
