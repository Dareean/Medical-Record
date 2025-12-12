import Logo from "../layout/Logo";

export default function Footer() {
  return (
    <footer className="flex items-center justify-between border-t border-gray-200 p-10">
        <Logo />

        <div className="flex items-center gap-5">
            <a href="" className="text-sm hover:text-blue-500">
                Github
            </a>

            <a href="" className="text-sm hover:text-blue-500">
                Discord
            </a>
        </div>
    </footer>
    
  )
}