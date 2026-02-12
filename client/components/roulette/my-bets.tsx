"use client";

import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useMyRouletteBets } from "@/hooks/use-roulette";
import { formatCurrency, formatDate } from "@/lib/utils";

const statusBadge = {
  pending: { label: "Pending", variant: "outline" as const },
  won: { label: "Won", variant: "default" as const },
  lost: { label: "Lost", variant: "destructive" as const },
} as const;

export function MyBets() {
  const { data: bets, isLoading } = useMyRouletteBets({ limit: 50 });

  if (isLoading) {
    return (
      <div className="text-sm text-muted-foreground p-4">Loading bets...</div>
    );
  }

  if (!bets || bets.length === 0) {
    return (
      <div className="text-sm text-muted-foreground p-4 text-center">
        No bets placed yet.
      </div>
    );
  }

  return (
    <div className="rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="text-xs">Type</TableHead>
            <TableHead className="text-xs">Value</TableHead>
            <TableHead className="text-xs">Amount</TableHead>
            <TableHead className="text-xs">Payout</TableHead>
            <TableHead className="text-xs">Status</TableHead>
            <TableHead className="text-xs">Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {bets.map((bet) => {
            const badge = statusBadge[bet.status];
            return (
              <TableRow key={bet.id}>
                <TableCell className="text-xs capitalize">
                  {bet.bet_type}
                </TableCell>
                <TableCell className="text-xs font-mono">
                  {bet.bet_value}
                </TableCell>
                <TableCell className="text-xs font-mono">
                  {formatCurrency(bet.amount)}
                </TableCell>
                <TableCell className="text-xs font-mono">
                  {formatCurrency(bet.payout)}
                </TableCell>
                <TableCell>
                  <Badge variant={badge.variant} className="text-[10px]">
                    {badge.label}
                  </Badge>
                </TableCell>
                <TableCell className="text-xs text-muted-foreground">
                  {formatDate(bet.created_at)}
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
}
