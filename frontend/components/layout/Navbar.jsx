import Logo from "../layout/Logo";

export default function Navbar() {
  return (
    <div className="fixed top-5 left-0 right-0 z-50 flex justify-center px-4">
      <nav className="flex items-center justify-between w-full max-w-4xl bg-white-400/90 backdrop-blur-md  px-6 py-3 rounded-full shadow-xl ring-1 ring-white/10">
        <div className="flex items-center gap-2">
          <div className="">
            <Logo />
          </div>
        </div>
        <div className="hidden md:flex items-center gap-8 font-medium text-sm">
          <a href="#hero" className="hover:text-red-400 transition-colors">
            Home
          </a>
          <a href="#services" className="hover:text-red-400 transition-colors">
            Service
          </a>
          <a href="#about" className="hover:text-red-400 transition-colors">
            About Us
          </a>
        </div>

        <div>
          <a
            href="/login"
            className="bg-white text-slate-800 px-5 py-2 rounded-full text-sm font-bold hover:bg-red-50 transition-colors cursor-pointer shadow-md"
          >
            Start Recording
          </a>
        </div>
      </nav>
    </div>
  );
}
