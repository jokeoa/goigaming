"use client";

import { X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { formatCurrency } from "@/lib/utils";
import { useRouletteStore } from "@/stores/roulette-store";

type PendingBetsProps = {
  readonly onPlaceAll: () => void;
  readonly isPlacing: boolean;
  readonly disabled: boolean;
};

export function PendingBets({
  onPlaceAll,
  isPlacing,
  disabled,
}: PendingBetsProps) {
  const pendingBets = useRouletteStore((s) => s.pendingBets);
  const removePendingBet = useRouletteStore((s) => s.removePendingBet);
  const clearPendingBets = useRouletteStore((s) => s.clearPendingBets);

  if (pendingBets.length === 0) {
    return null;
  }

  const total = pendingBets.reduce(
    (sum, bet) => sum + Number.parseFloat(bet.amount),
    0,
  );

  return (
    <div className="space-y-3 rounded-lg border border-border bg-card p-4">
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium">
          Pending Bets ({pendingBets.length})
        </h3>
        <span className="text-sm font-mono font-medium">
          Total: {formatCurrency(total)}
        </span>
      </div>

      <div className="space-y-1.5">
        {pendingBets.map((bet) => (
          <div
            key={bet.id}
            className="flex items-center justify-between rounded border border-border bg-background px-3 py-1.5 text-xs"
          >
            <span>
              {bet.bet_type}
              {bet.bet_value !== bet.bet_type ? ` (${bet.bet_value})` : ""} -{" "}
              {formatCurrency(bet.amount)}
            </span>
            <button
              type="button"
              onClick={() => removePendingBet(bet.id)}
              className="text-muted-foreground hover:text-destructive"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </div>
        ))}
      </div>

      <div className="flex gap-2">
        <Button
          size="sm"
          onClick={onPlaceAll}
          disabled={isPlacing || disabled}
          className="flex-1"
        >
          {isPlacing ? "Placing..." : "Place All Bets"}
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={clearPendingBets}
          disabled={isPlacing}
        >
          Clear All
        </Button>
      </div>
    </div>
  );
}
