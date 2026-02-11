import Link from "next/link";
import { Button } from "@/components/ui/button";

export function CtaSection() {
  return (
    <section className="border-t border-border py-16">
      <div className="mx-auto flex max-w-2xl flex-col items-center px-4 text-center">
        <h2 className="text-2xl font-semibold tracking-tight">
          Ready to play?
        </h2>
        <p className="mt-2 text-muted-foreground">
          Create your free account and join a table in under a minute.
        </p>
        <Button asChild className="mt-6" size="lg">
          <Link href="/register">Create free account</Link>
        </Button>
      </div>
    </section>
  );
}
