"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Menu, X, Github } from "lucide-react";

const navLinks = [
  { label: "Features", href: "#features" },
  { label: "Use Cases", href: "#use-cases" },
  { label: "Architecture", href: "#architecture" },
  { label: "Open Source", href: "#open-source" },
];

export function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const [mobileOpen, setMobileOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 50);
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <motion.nav
      initial={{ y: -100 }}
      animate={{ y: 0 }}
      transition={{ duration: 0.5 }}
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled
          ? "bg-[#0A0B14]/80 backdrop-blur-xl border-b border-[#1E293B]"
          : "bg-transparent"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <a href="#" className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-[#6366F1] to-[#8B5CF6] flex items-center justify-center">
              <span className="text-white font-bold text-sm">K8</span>
            </div>
            <span className="font-semibold text-lg text-white">
              clawdlinux
            </span>
          </a>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-8">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                className="text-sm text-[#94A3B8] hover:text-white transition-colors"
              >
                {link.label}
              </a>
            ))}
          </div>

          {/* CTA Buttons */}
          <div className="hidden md:flex items-center gap-3">
            <a
              href="https://github.com/shreyanshjain7174/agentic-k8s-operator"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-4 py-2 text-sm text-[#94A3B8] hover:text-white transition-colors"
            >
              <Github size={18} />
              GitHub
            </a>
            <a
              href="#waitlist"
              className="px-5 py-2 text-sm font-medium rounded-full bg-gradient-to-r from-[#6366F1] to-[#8B5CF6] text-white hover:shadow-lg hover:shadow-[#6366F1]/25 transition-all"
            >
              Join Waitlist
            </a>
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden text-white"
            onClick={() => setMobileOpen(!mobileOpen)}
          >
            {mobileOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>
      </div>

      {/* Mobile menu */}
      {mobileOpen && (
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          className="md:hidden bg-[#0A0B14]/95 backdrop-blur-xl border-b border-[#1E293B]"
        >
          <div className="px-4 py-4 space-y-3">
            {navLinks.map((link) => (
              <a
                key={link.href}
                href={link.href}
                onClick={() => setMobileOpen(false)}
                className="block text-sm text-[#94A3B8] hover:text-white transition-colors py-2"
              >
                {link.label}
              </a>
            ))}
            <a
              href="#waitlist"
              onClick={() => setMobileOpen(false)}
              className="block text-center px-5 py-2 text-sm font-medium rounded-full bg-gradient-to-r from-[#6366F1] to-[#8B5CF6] text-white"
            >
              Join Waitlist
            </a>
          </div>
        </motion.div>
      )}
    </motion.nav>
  );
}
