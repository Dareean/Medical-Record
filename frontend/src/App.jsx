import { ArrowRight } from "lucide-react";
import Navbar from "../components/layout/Navbar";
import Footer from "../components/layout/Footer";

function App() {
  return (
    <div className="">
      <Navbar />

      <main className="py-10">
        <section
          id="hero"
          className="flex items-center justify-between gap-2 px-10 py-20"
        >
          <div className="space-y-6 max-w-xl">
            <h1 className="text-7xl font-black font-serif">
              Build with <span className="text-blue-400">precision.</span>
            </h1>
            <h2 className="text-gray-500 text-lg">
              A minimalist toolkit for developers who value clean code and
              simple workflows. No bloat. Just tools.
            </h2>
            <div className="flex items-center gap-4 text-sm">
              <button className="group flex items-center gap-2 text-white bg-blue-400 px-6 py-2 rounded-md cursor-pointer hover:bg-blue-300">
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

          <div className="border border-gray-200 rounded-lg p-6 space-y-2 w-full max-w-2xl">
            <div className="flex items-center gap-2">
              <div className="w-3 h-3 bg-red-400 rounded-full"></div>
              <div className="w-3 h-3 bg-yellow-400 rounded-full"></div>
              <div className="w-3 h-3 bg-green-400 rounded-full"></div>
            </div>

            <pre className="text-gray-500">{`> hammercode init

Creating project...
✓ Dependencies installed
✓ Config generated
✓ Ready to build

> hammercode dev
Server running at :3000`}</pre>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}

export default App;
