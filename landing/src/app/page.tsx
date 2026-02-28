import { Navbar } from "@/components/sections/Navbar";
import { Hero } from "@/components/sections/Hero";
import { MetricsBar } from "@/components/sections/MetricsBar";
import { Features } from "@/components/sections/Features";
import { UseCases } from "@/components/sections/UseCases";
import { Architecture } from "@/components/sections/Architecture";
import { OpenSource } from "@/components/sections/OpenSource";
import { Waitlist } from "@/components/sections/Waitlist";
import { Footer } from "@/components/sections/Footer";

export default function Home() {
  return (
    <>
      <Navbar />
      <main>
        <Hero />
        <MetricsBar />
        <Features />
        <UseCases />
        <Architecture />
        <OpenSource />
        <Waitlist />
      </main>
      <Footer />
    </>
  );
}
