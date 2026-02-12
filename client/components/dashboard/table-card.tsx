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
  return (
    <Card className="border-border">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm">{table.name}</CardTitle>
        <Badge
          variant={table.status === "active" ? "default" : "secondary"}
          className="text-xs"
        >
          {table.status === "waiting"
            ? "Open"
            : table.status === "active"
              ? "In Progress"
              : "Closed"}
        </Badge>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>Max: {table.max_players} players</span>
          <span>
            Blinds: {formatCurrency(table.small_blind)}/
            {formatCurrency(table.big_blind)}
          </span>
        </div>
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>
            Buy-in: {formatCurrency(table.min_buy_in)} -{" "}
            {formatCurrency(table.max_buy_in)}
          </span>
        </div>
        <Button
          asChild
          size="sm"
          className="w-full"
          disabled={table.status === "closed"}
        >
          <Link href={`/table/${table.id}`}>
            {table.status === "active" ? "Join Table" : "Open Table"}
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
