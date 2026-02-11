import { ArrowRight } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export function HeroSection() {
  return (
    <section className="flex flex-col items-center justify-center px-4 py-24 text-center sm:py-32">
      <div className="mb-4 inline-flex items-center rounded-full border border-border bg-card px-3 py-1 text-xs text-muted-foreground">
        Real-time multiplayer poker
      </div>
      <h1 className="max-w-2xl text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
        Play poker with
        <span className="text-primary"> real players</span>
      </h1>
      <p className="mt-6 max-w-lg text-lg text-muted-foreground">
        Join tables, manage your wallet, and compete in real-time Texas
        Hold&apos;em games on a modern, sleek platform.
      </p>
      <div className="mt-8 flex gap-4">
        <Button asChild size="lg">
          <Link href="/register">
            Get started
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
        <Button asChild variant="outline" size="lg">
          <Link href="/login">Sign in</Link>
        </Button>
      </div>
    </section>
  );
}
