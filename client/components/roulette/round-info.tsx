"use client";

import { Badge } from "@/components/ui/badge";
import type { RouletteRound } from "@/types/roulette";

type RoundInfoProps = {
  readonly round: RouletteRound | undefined;
  readonly isLoading: boolean;
};

function getRoundStatus(round: RouletteRound): {
  label: string;
  variant: "default" | "secondary" | "destructive" | "outline";
} {
  if (round.settled_at) {
    return { label: "Settled", variant: "secondary" };
  }
  if (round.result !== null) {
    return { label: "Spinning", variant: "destructive" };
  }
  return { label: "Betting Open", variant: "default" };
}

function BettingCountdown({ endsAt }: { readonly endsAt: string }) {
  const endTime = new Date(endsAt).getTime();
  const now = Date.now();
  const remaining = Math.max(0, Math.ceil((endTime - now) / 1000));

  if (remaining <= 0) {
    return <span className="text-xs text-destructive">Betting closed</span>;
  }

  return (
    <span className="text-xs font-mono text-muted-foreground">
      {remaining}s remaining
    </span>
  );
}

export function RoundInfo({ round, isLoading }: RoundInfoProps) {
  if (isLoading) {
    return (
      <div className="flex items-center justify-between rounded-lg border border-border bg-card p-3">
        <span className="text-sm text-muted-foreground">Loading round...</span>
      </div>
    );
  }

  if (!round) {
    return (
      <div className="flex items-center justify-center rounded-lg border border-dashed border-border bg-card/50 p-4 text-sm text-muted-foreground">
        No active round. Waiting for next round to start...
      </div>
    );
  }

  const status = getRoundStatus(round);

  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border bg-card p-3 text-sm">
      <div className="flex items-center gap-3">
        <span className="font-medium">Round #{round.round_number}</span>
        <Badge variant={status.variant} className="text-xs">
          {status.label}
        </Badge>
      </div>
      {round.betting_ends_at && !round.settled_at && (
        <BettingCountdown endsAt={round.betting_ends_at} />
      )}
    </div>
  );
}
