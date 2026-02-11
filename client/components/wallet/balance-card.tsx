"use client";

import { Wallet } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useWalletBalance } from "@/hooks/use-wallet";
import { formatCurrency } from "@/lib/utils";

export function BalanceCard() {
  const { data: wallet, isLoading } = useWalletBalance();

  return (
    <Card className="border-border">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          Current Balance
        </CardTitle>
        <Wallet className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Skeleton className="h-8 w-32" />
        ) : (
          <p className="text-3xl font-bold font-mono text-primary">
            {wallet ? formatCurrency(wallet.balance) : "$0.00"}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
