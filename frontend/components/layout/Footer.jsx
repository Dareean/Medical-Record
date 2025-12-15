import {
  Facebook,
  Twitter,
  Instagram,
  Linkedin,
  Mail,
  MapPin,
  Phone,
} from "lucide-react";
import Logo from "../layout/Logo";

export default function Footer() {
  return (
    <footer className="bg-white text-slate-300 pt-16 pb-8 shadow-2xl">
      <div className="max-w-7xl mx-auto px-6 md:px-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5 mb-12">
          <div className="space-y-4">
            <div className="scale-100 text-slate-500">
              <Logo />
            </div>
            <p className="text-slate-400 text-sm pr-4">
              A trusted digital medical record platform. Manage patient data
              securely, quickly, and efficiently for better healthcare services.
            </p>
            <div className="flex gap-4 pt-2">
              <SocialIcon icon={<Facebook size={18} />} />
              <SocialIcon icon={<Twitter size={18} />} />
              <SocialIcon icon={<Instagram size={18} />} />
              <SocialIcon icon={<Linkedin size={18} />} />
            </div>
          </div>

          <div>
            <h3 className="text-slate-500 font-bold mb-6 ">Quick Links</h3>
            <ul className="space-y-3 text-sm text-slate-400">
              <li>
                <FooterLink href="#hero">Home</FooterLink>
              </li>
              <li>
                <FooterLink href="#services">Our Services</FooterLink>
              </li>
              <li>
                <FooterLink href="#about">About Us</FooterLink>
              </li>
            </ul>
          </div>

          <div>
            <h3 className="text-slate-500 font-bold mb-6">Contact Us</h3>
            <ul className="space-y-4 text-sm">
              <li className="flex items-start gap-3">
                <MapPin className="w-5 h-5 text-red-400 shrink-0" />
                <span>Jl. Veteran No. 108A Palu. Indonesia</span>
              </li>
              <li className="flex items-center gap-3">
                <Mail className="w-5 h-5 text-red-400 shrink-0" />
                <span>support@medicflow.id</span>
              </li>
              <li className="flex items-center gap-3">
                <Phone className="w-5 h-5 text-red-400 shrink-0" />
                <span>+62 12 3456 7890</span>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-white-800 pt-8 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-slate-500 text-xs">
            © {new Date().getFullYear()} MedicFlow. All rights reserved.
          </p>
          <div className="flex gap-6 text-xs text-slate-500">
            <a href="#" className="hover:text-red-400 transition-colors">
              Privacy
            </a>
            <a href="#" className="hover:text-red-400 transition-colors">
              Cookies
            </a>
            <a href="#" className="hover:text-red-400 transition-colors">
              Terms
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}

function FooterLink({ href, children }) {
  return (
    <a
      href={href}
      className="hover:text-red-400 hover:translate-x-1 transition-all inline-block"
    >
      {children}
    </a>
  );
}

function SocialIcon({ icon }) {
  return (
    <a
      href="#"
      className="w-8 h-8 rounded-full bg-white flex items-center justify-center text-black-400 hover:bg-red-500 hover:text-white transition-all duration-300"
    >
      {icon}
    </a>
  );
}
