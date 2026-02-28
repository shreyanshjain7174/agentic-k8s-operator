import { Github } from "lucide-react";

const footerLinks = {
  Product: [
    { label: "Features", href: "#features" },
    { label: "Use Cases", href: "#use-cases" },
    { label: "Architecture", href: "#architecture" },
    { label: "Documentation", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator#readme" },
  ],
  Resources: [
    { label: "GitHub", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator" },
    { label: "Helm Charts", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator/tree/main/charts" },
    { label: "API Reference", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator/tree/main/api" },
    { label: "License", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator/blob/main/LICENSE" },
  ],
  Community: [
    { label: "Open Source", href: "#open-source" },
    { label: "Contributing", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator" },
    { label: "Issues", href: "https://github.com/shreyanshjain7174/agentic-k8s-operator/issues" },
  ],
};

export function Footer() {
  return (
    <footer className="border-t border-[#1E293B] bg-[#0A0B14]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <div className="flex items-center gap-2 mb-4">
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-[#6366F1] to-[#8B5CF6] flex items-center justify-center">
                <span className="text-white font-bold text-sm">K8</span>
              </div>
              <span className="font-semibold text-white">clawdlinux</span>
            </div>
            <p className="text-sm text-[#64748B] mb-4 max-w-xs">
              Production-grade Kubernetes operator for orchestrating autonomous
              AI agent workloads.
            </p>
            <a
              href="https://github.com/shreyanshjain7174/agentic-k8s-operator"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[#94A3B8] hover:text-white transition-colors"
            >
              <Github size={20} />
            </a>
          </div>

          {/* Link columns */}
          {Object.entries(footerLinks).map(([title, links]) => (
            <div key={title}>
              <h4 className="text-sm font-semibold text-white mb-4">{title}</h4>
              <ul className="space-y-2">
                {links.map((link) => (
                  <li key={link.label}>
                    <a
                      href={link.href}
                      target={link.href.startsWith("http") ? "_blank" : undefined}
                      rel={link.href.startsWith("http") ? "noopener noreferrer" : undefined}
                      className="text-sm text-[#64748B] hover:text-[#94A3B8] transition-colors"
                    >
                      {link.label}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom bar */}
        <div className="mt-12 pt-8 border-t border-[#1E293B] flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-xs text-[#64748B]">
            &copy; {new Date().getFullYear()} clawdlinux.org. Apache 2.0 Licensed.
          </p>
          <p className="text-xs text-[#64748B]">
            Built with care for the Kubernetes community
          </p>
        </div>
      </div>
    </footer>
  );
}
