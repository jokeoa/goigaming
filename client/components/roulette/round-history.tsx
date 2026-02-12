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
import { useRoundHistory } from "@/hooks/use-roulette";
import { formatDate } from "@/lib/utils";

type RoundHistoryProps = {
  readonly tableId: string;
};

const colorBadgeVariant = {
  red: "destructive",
  black: "secondary",
  green: "default",
} as const;

export function RoundHistory({ tableId }: RoundHistoryProps) {
  const { data: rounds, isLoading } = useRoundHistory(tableId, { limit: 20 });

  if (isLoading) {
    return (
      <div className="text-sm text-muted-foreground p-4">
        Loading history...
      </div>
    );
  }

  if (!rounds || rounds.length === 0) {
    return (
      <div className="text-sm text-muted-foreground p-4 text-center">
        No round history yet.
      </div>
    );
  }

  return (
    <div className="rounded-lg border border-border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="text-xs">Round</TableHead>
            <TableHead className="text-xs">Result</TableHead>
            <TableHead className="text-xs">Color</TableHead>
            <TableHead className="text-xs">Settled</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rounds.map((round) => (
            <TableRow key={round.id}>
              <TableCell className="text-xs font-mono">
                #{round.round_number}
              </TableCell>
              <TableCell className="text-xs font-mono font-medium">
                {round.result ?? "-"}
              </TableCell>
              <TableCell>
                {round.result_color && (
                  <Badge
                    variant={
                      colorBadgeVariant[
                        round.result_color as keyof typeof colorBadgeVariant
                      ] ?? "outline"
                    }
                    className="text-[10px]"
                  >
                    {round.result_color}
                  </Badge>
                )}
              </TableCell>
              <TableCell className="text-xs text-muted-foreground">
                {round.settled_at ? formatDate(round.settled_at) : "-"}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
