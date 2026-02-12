"use client";

import { toast } from "sonner";
import { BetSelector } from "@/components/roulette/bet-selector";
import { PendingBets } from "@/components/roulette/pending-bets";
import { RoundHistory } from "@/components/roulette/round-history";
import { RoundInfo } from "@/components/roulette/round-info";
import { RoundResult } from "@/components/roulette/round-result";
import { useCurrentRound, usePlaceBet } from "@/hooks/use-roulette";
import { useRouletteStore } from "@/stores/roulette-store";

type RouletteGameViewProps = {
  readonly tableId: string;
};

export function RouletteGameView({ tableId }: RouletteGameViewProps) {
  const { data: round, isLoading } = useCurrentRound(tableId);
  const placeBetMutation = usePlaceBet();
  const pendingBets = useRouletteStore((s) => s.pendingBets);
  const clearPendingBets = useRouletteStore((s) => s.clearPendingBets);

  const isBettingOpen = round && !round.settled_at && round.result === null;

  const handlePlaceAll = async () => {
    if (!round || pendingBets.length === 0) return;

    try {
      for (const bet of pendingBets) {
        await placeBetMutation.mutateAsync({
          tableId,
          body: {
            round_id: round.id,
            bet_type: bet.bet_type,
            bet_value: bet.bet_value,
            amount: bet.amount,
          },
        });
      }
      clearPendingBets();
      toast.success("All bets placed successfully");
    } catch (error) {
      const message =
        error instanceof Error ? error.message : "Failed to place bets";
      toast.error(message);
    }
  };

  return (
    <div className="space-y-6">
      <RoundInfo round={round} isLoading={isLoading} />

      <RoundResult round={round} />

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="space-y-4">
          <h2 className="text-sm font-semibold">Place Bets</h2>
          <BetSelector />
          <PendingBets
            onPlaceAll={handlePlaceAll}
            isPlacing={placeBetMutation.isPending}
            disabled={!isBettingOpen}
          />
        </div>

        <div className="space-y-4">
          <h2 className="text-sm font-semibold">Round History</h2>
          <RoundHistory tableId={tableId} />
        </div>
      </div>
    </div>
  );
}
