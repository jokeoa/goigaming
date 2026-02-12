import { Badge } from "@/components/ui/badge";
import { formatCurrency } from "@/lib/utils";
import type { GameStage, PokerTable } from "@/types/game";

type TableInfoProps = {
  readonly table: PokerTable;
  readonly stage: GameStage;
  readonly isConnected: boolean;
};

const stageLabels: Record<GameStage, string> = {
  waiting: "Waiting",
  preflop: "Pre-Flop",
  flop: "Flop",
  turn: "Turn",
  river: "River",
  showdown: "Showdown",
} as const;

export function TableInfo({ table, stage, isConnected }: TableInfoProps) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 rounded-lg border border-border bg-card p-3 text-sm">
      <div className="flex items-center gap-3">
        <span className="font-medium">{table.name}</span>
        <Badge variant="outline" className="text-xs">
          {stageLabels[stage]}
        </Badge>
      </div>
      <div className="flex items-center gap-4 text-xs text-muted-foreground">
        <span>
          Blinds: {formatCurrency(table.smallBlind)}/
          {formatCurrency(table.bigBlind)}
        </span>
        <span>
          Players: {table.currentPlayers}/{table.maxPlayers}
        </span>
        <div className="flex items-center gap-1">
          <div
            className={`h-2 w-2 rounded-full ${isConnected ? "bg-primary" : "bg-destructive"}`}
          />
          <span>{isConnected ? "Connected" : "Disconnected"}</span>
        </div>
      </div>
    </div>
  );
}
