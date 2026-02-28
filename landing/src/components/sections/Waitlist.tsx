"use client";

import { useState, FormEvent } from "react";
import { AnimatedSection } from "../ui/AnimatedSection";
import { Send, CheckCircle, AlertCircle, Loader2 } from "lucide-react";

export function Waitlist() {
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [status, setStatus] = useState<"idle" | "loading" | "success" | "error">("idle");

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!email) return;

    setStatus("loading");

    try {
      const res = await fetch("/api/waitlist", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, name }),
      });

      if (res.ok) {
        setStatus("success");
        setEmail("");
        setName("");
      } else {
        setStatus("error");
      }
    } catch {
      setStatus("error");
    }
  };

  return (
    <section id="waitlist" className="py-24 relative">
      {/* Background glow */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[400px] bg-[#6366F1]/5 rounded-full blur-[120px]" />
      </div>

      <div className="relative max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        <AnimatedSection>
          <div className="glass-card p-8 sm:p-12 text-center glow-effect">
            <h2 className="text-3xl sm:text-4xl font-bold mb-4">
              Be the First to Deploy{" "}
              <span className="gradient-text">AI Agents</span> in Your Cluster
            </h2>
            <p className="text-[#94A3B8] mb-8 max-w-lg mx-auto">
              Join the waitlist for early access to the commercial tier and
              design partner program. Get priority support and help shape the
              roadmap.
            </p>

            {status === "success" ? (
              <div className="flex flex-col items-center gap-3 py-4">
                <CheckCircle size={48} className="text-[#10B981]" />
                <p className="text-lg font-medium text-[#10B981]">
                  You&apos;re on the list!
                </p>
                <p className="text-sm text-[#94A3B8]">
                  We&apos;ll be in touch soon with early access details.
                </p>
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="grid sm:grid-cols-2 gap-4">
                  <input
                    type="text"
                    placeholder="Your name"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    className="w-full px-4 py-3 rounded-xl bg-[#0A0B14] border border-[#1E293B] text-white placeholder-[#64748B] focus:outline-none focus:border-[#6366F1] transition-colors"
                  />
                  <input
                    type="email"
                    placeholder="Your email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    className="w-full px-4 py-3 rounded-xl bg-[#0A0B14] border border-[#1E293B] text-white placeholder-[#64748B] focus:outline-none focus:border-[#6366F1] transition-colors"
                  />
                </div>

                <button
                  type="submit"
                  disabled={status === "loading"}
                  className="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-8 py-3 rounded-full bg-gradient-to-r from-[#6366F1] to-[#8B5CF6] text-white font-medium hover:shadow-lg hover:shadow-[#6366F1]/25 transition-all disabled:opacity-50"
                >
                  {status === "loading" ? (
                    <>
                      <Loader2 size={18} className="animate-spin" />
                      Joining...
                    </>
                  ) : (
                    <>
                      <Send size={18} />
                      Join Waitlist
                    </>
                  )}
                </button>

                {status === "error" && (
                  <div className="flex items-center justify-center gap-2 text-sm text-[#EF4444]">
                    <AlertCircle size={16} />
                    Something went wrong. Please try again.
                  </div>
                )}

                <p className="text-xs text-[#64748B]">
                  No spam. Unsubscribe anytime.
                </p>
              </form>
            )}
          </div>
        </AnimatedSection>
      </div>
    </section>
  );
}
