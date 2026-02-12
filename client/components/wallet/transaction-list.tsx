"use client";

import { Receipt } from "lucide-react";
import { EmptyState } from "@/components/shared/empty-state";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useTransactions } from "@/hooks/use-transactions";
import { formatCurrency, formatDate } from "@/lib/utils";

export function TransactionList() {
  const { data: transactions, isLoading } = useTransactions({ limit: 20 });

  if (isLoading) {
    return (
      <Card className="border-border">
        <CardHeader>
          <CardTitle className="text-sm">Recent Transactions</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={`tx-skeleton-${i}`} className="h-10 w-full" />
          ))}
        </CardContent>
      </Card>
    );
  }

  if (!transactions || transactions.length === 0) {
    return (
      <EmptyState
        icon={<Receipt className="h-8 w-8" />}
        title="No transactions yet"
        description="Your transaction history will appear here."
      />
    );
  }

  return (
    <Card className="border-border">
      <CardHeader>
        <CardTitle className="text-sm">Recent Transactions</CardTitle>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Type</TableHead>
              <TableHead>Amount</TableHead>
              <TableHead>Balance After</TableHead>
              <TableHead className="hidden sm:table-cell">Date</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {transactions.map((tx) => {
              const isPositive = Number.parseFloat(tx.amount) > 0;
              return (
                <TableRow key={tx.id}>
                  <TableCell>
                    <Badge
                      variant={isPositive ? "default" : "secondary"}
                      className="text-xs"
                    >
                      {tx.reference_type}
                    </Badge>
                  </TableCell>
                  <TableCell
                    className={`font-mono text-sm ${isPositive ? "text-primary" : "text-destructive"}`}
                  >
                    {isPositive ? "+" : ""}
                    {formatCurrency(tx.amount)}
                  </TableCell>
                  <TableCell className="font-mono text-sm">
                    {formatCurrency(tx.balance_after)}
                  </TableCell>
                  <TableCell className="hidden text-xs text-muted-foreground sm:table-cell">
                    {formatDate(tx.created_at)}
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
