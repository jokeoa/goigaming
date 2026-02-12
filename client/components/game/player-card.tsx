import { Badge } from "@/components/ui/badge";
import { cn, formatCurrency } from "@/lib/utils";
import type { WSPlayerInfo } from "@/types/game";

type PlayerCardProps = {
  readonly player: WSPlayerInfo;
  readonly isCurrentTurn: boolean;
};

export function PlayerCard({ player, isCurrentTurn }: PlayerCardProps) {
  const isFolded = player.status === "folded";

  return (
    <div
      className={cn(
        "flex flex-col items-center gap-1 rounded-lg border bg-card p-2 text-center transition-colors",
        isCurrentTurn && "border-primary ring-1 ring-primary",
        isFolded && "opacity-50",
        !isCurrentTurn && "border-border",
      )}
    >
      <div className="flex items-center gap-1">
        <span className="text-xs font-medium truncate max-w-[72px]">
          {player.username}
        </span>
        {player.is_dealer && (
          <Badge
            variant="outline"
            className="h-4 px-1 text-[10px] text-gold border-gold"
          >
            D
          </Badge>
        )}
      </div>
      <span className="text-xs font-mono text-muted-foreground">
        {formatCurrency(player.stack)}
      </span>
      {Number.parseFloat(player.bet_amount) > 0 && (
        <span className="text-[10px] font-mono text-primary">
          Bet: {formatCurrency(player.bet_amount)}
        </span>
      )}
      {isFolded && <span className="text-[10px] text-destructive">Folded</span>}
    </div>
  );
}
