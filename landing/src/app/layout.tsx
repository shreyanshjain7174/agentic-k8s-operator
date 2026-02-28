import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Agentic K8s Operator | Deploy AI Agents Inside Your Kubernetes Cluster",
  description:
    "Production-grade Kubernetes operator for orchestrating autonomous AI agent workloads. One Helm command. Zero data leaves your infrastructure. Apache 2.0 open source.",
  keywords: [
    "kubernetes",
    "ai agents",
    "k8s operator",
    "multi-agent orchestration",
    "MCP",
    "LangGraph",
    "open source",
  ],
  openGraph: {
    title: "Agentic K8s Operator",
    description:
      "Deploy AI Agents Inside Your Kubernetes Cluster. Production-grade. One Helm command.",
    url: "https://clawdlinux.org",
    siteName: "Agentic K8s Operator",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Agentic K8s Operator",
    description:
      "Deploy AI Agents Inside Your Kubernetes Cluster. Production-grade. One Helm command.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className="dark">
      <body className="antialiased bg-[#0A0B14] text-[#F8FAFC] font-sans">
        {children}
      </body>
    </html>
  );
}
