import { Gamepad2, Shield, Wallet, Zap } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const features = [
  {
    icon: Gamepad2,
    title: "Live Tables",
    description:
      "Join real-time poker tables with players from around the world.",
  },
  {
    icon: Wallet,
    title: "Instant Wallet",
    description:
      "Deposit and withdraw instantly. Track every transaction in real time.",
  },
  {
    icon: Zap,
    title: "Fast Gameplay",
    description:
      "WebSocket-powered game engine for zero-latency card dealing and actions.",
  },
  {
    icon: Shield,
    title: "Secure Platform",
    description:
      "JWT authentication, encrypted connections, and fair game logic.",
  },
] as const;

export function FeaturesSection() {
  return (
    <section className="mx-auto max-w-5xl px-4 py-16">
      <h2 className="mb-8 text-center text-2xl font-semibold tracking-tight">
        Built for serious players
      </h2>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {features.map((feature) => (
          <Card key={feature.title} className="border-border bg-card">
            <CardHeader className="pb-2">
              <feature.icon className="mb-2 h-5 w-5 text-primary" />
              <CardTitle className="text-sm">{feature.title}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-xs text-muted-foreground">
                {feature.description}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </section>
  );
}
