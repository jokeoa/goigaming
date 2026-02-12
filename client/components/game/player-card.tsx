import { Badge } from "@/components/ui/badge";
import { cn, formatCurrency } from "@/lib/utils";
import type { PlayerInfo } from "@/types/game";

type PlayerCardProps = {
  readonly player: PlayerInfo;
  readonly isCurrentTurn: boolean;
};

export function PlayerCard({ player, isCurrentTurn }: PlayerCardProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center gap-1 rounded-lg border bg-card p-2 text-center transition-colors",
        isCurrentTurn && "border-primary ring-1 ring-primary",
        player.isFolded && "opacity-50",
        !isCurrentTurn && "border-border",
      )}
    >
      <div className="flex items-center gap-1">
        <span className="text-xs font-medium truncate max-w-[72px]">
          {player.username}
        </span>
        {player.isDealer && (
          <Badge
            variant="outline"
            className="h-4 px-1 text-[10px] text-gold border-gold"
          >
            D
          </Badge>
        )}
      </div>
      <span className="text-xs font-mono text-muted-foreground">
        {formatCurrency(player.chips)}
      </span>
      {player.currentBet > 0 && (
        <span className="text-[10px] font-mono text-primary">
          Bet: {formatCurrency(player.currentBet)}
        </span>
      )}
      {player.isFolded && (
        <span className="text-[10px] text-destructive">Folded</span>
      )}
    </div>
  );
}
