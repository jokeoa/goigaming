import Link from "next/link";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { formatCurrency } from "@/lib/utils";
import type { RouletteTable } from "@/types/roulette";

type RouletteTableCardProps = {
  readonly table: RouletteTable;
};

const statusVariant = {
  active: "default",
  inactive: "secondary",
  maintenance: "destructive",
} as const;

export function RouletteTableCard({ table }: RouletteTableCardProps) {
  return (
    <Card className="border-border">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm">{table.name}</CardTitle>
        <Badge variant={statusVariant[table.status]} className="text-xs">
          {table.status}
        </Badge>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="flex items-center justify-between text-xs text-muted-foreground">
          <span>Min: {formatCurrency(table.min_bet)}</span>
          <span>Max: {formatCurrency(table.max_bet)}</span>
        </div>
        <Button
          asChild
          size="sm"
          className="w-full"
          disabled={table.status !== "active"}
        >
          <Link href={`/roulette/${table.id}`}>Play</Link>
        </Button>
      </CardContent>
    </Card>
  );
}
