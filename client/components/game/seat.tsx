import { PlayerCard } from "@/components/game/player-card";
import { cn } from "@/lib/utils";
import type { PlayerInfo } from "@/types/game";

type SeatProps = {
  readonly player?: PlayerInfo;
  readonly position: number;
  readonly isCurrentTurn: boolean;
  readonly className?: string;
};

export function Seat({
  player,
  position,
  isCurrentTurn,
  className,
}: SeatProps) {
  if (!player) {
    return (
      <div
        className={cn(
          "flex h-20 w-24 items-center justify-center rounded-lg border border-dashed border-border bg-card/30 text-xs text-muted-foreground",
          className,
        )}
      >
        Seat {position + 1}
      </div>
    );
  }

  return (
    <div className={className}>
      <PlayerCard player={player} isCurrentTurn={isCurrentTurn} />
    </div>
  );
}
