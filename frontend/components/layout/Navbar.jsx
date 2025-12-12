import Logo from "../layout/Logo";

export default function Navbar() {
  return (
    <header className="fixed w-full top-0 flex items-center justify-between bg-white border-b border-gray-200 px-10 py-5">
      <Logo />

      <div className="flex items-center gap-2 font-medium">
        <a href="" className="text-sm hover:text-blue-500">
          Home
        </a>

        <a href="" className="text-sm hover:text-blue-500">
          Tentang
        </a>
      </div>
    </header>
  );
}
