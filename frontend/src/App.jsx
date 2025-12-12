import { ArrowRight } from "lucide-react";
import Navbar from "../components/layout/Navbar";
import Footer from "../components/layout/Footer";
import { useRef } from "React";

function App() {
  const homeRef = useRef();
  return (
    <div className="">
      <Navbar />

      <main className="py-10">
        <section
          id="hero"
          className="flex items-center justify-between gap-2 px-25 py-10"
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
              <button className="group flex items-center gap-2 text-white bg-red-400 px-6 py-2 rounded-md cursor-pointer hover:bg-blue-300">
                Get Started
                <ArrowRight
                  size={16}
                  className="group-hover:translate-x-1 transition-transform"
                />
              </button>
              <button className="cursor-pointer rounded-md border border-gray-400 px-4 py-2 hover:bg-blue-400 transition-colors duration-300 hover:text-white">
                View Docs
              </button>
            </div>
          </div>

          <div className="p-6 space-y-2 w-full max-w-2xl">
            <img
              className="w-full h-auto object-cover"
              src="/dashboard-medical-record.png"
              alt="Medical Illustration"
            />
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}

export default App;
