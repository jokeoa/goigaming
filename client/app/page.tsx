import Link from "next/link";
import { CtaSection } from "@/components/landing/cta-section";
import { FeaturesSection } from "@/components/landing/features-section";
import { HeroSection } from "@/components/landing/hero-section";

export default function LandingPage() {
  return (
    <div className="flex min-h-screen flex-col bg-background">
      <header className="flex h-14 items-center justify-between border-b border-border px-4 lg:px-6">
        <Link href="/" className="flex items-center gap-2">
          <img
            src="/logo.svg"
            alt="GoiGaming logo"
            className="h-7 w-auto shrink-0"
          />
          <span className="font-mono text-sm font-semibold">GoiGaming</span>
        </Link>
        <nav className="flex items-center gap-4 text-sm">
          <Link
            href="/login"
            className="text-muted-foreground transition-colors hover:text-foreground"
          >
            Sign in
          </Link>
          <Link
            href="/register"
            className="rounded-lg bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
          >
            Sign up
          </Link>
        </nav>
      </header>
      <main className="flex-1">
        <HeroSection />
        <FeaturesSection />
        <CtaSection />
      </main>
      <footer className="border-t border-border py-6 text-center text-xs text-muted-foreground">
        GoiGaming &mdash; Play responsibly
      </footer>
    </div>
  );
}
