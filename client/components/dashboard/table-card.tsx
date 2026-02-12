import { Users } from "lucide-react";
import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatCurrency } from "@/lib/utils";
import type { PokerTable } from "@/types/game";

type TableCardProps = {
  readonly table: PokerTable;
};

export function TableCard({ table }: TableCardProps) {
  const isFull = table.currentPlayers >= table.maxPlayers;

  return (
    <Card className="border-border">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm">{table.name}</CardTitle>
        <Badge variant={isFull ? "secondary" : "default"} className="text-xs">
          {isFull ? "Full" : table.stage === "waiting" ? "Open" : "In Progress"}
        </Badge>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <div className="flex items-center gap-1">
            <Users className="h-3 w-3" />
            <span>
              {table.currentPlayers}/{table.maxPlayers}
            </span>
          </div>
          <span>
            Blinds: {formatCurrency(table.smallBlind)}/
            {formatCurrency(table.bigBlind)}
          </span>
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>
            Buy-in: {formatCurrency(table.minBuyIn)} -{" "}
            {formatCurrency(table.maxBuyIn)}
          </span>
        </div>
        <Button asChild size="sm" className="w-full" disabled={isFull}>
          <Link href={`/table/${table.id}`}>
            {isFull ? "Spectate" : "Join Table"}
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
