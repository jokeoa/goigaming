"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useWalletBalance } from "@/hooks/use-wallet";
import { formatCurrency } from "@/lib/utils";

export function ProfileStats() {
  const { data: wallet, isLoading } = useWalletBalance();

  const stats = [
    {
      label: "Balance",
      value: wallet ? formatCurrency(wallet.balance) : "$0.00",
    },
    { label: "Games Played", value: "0" },
    { label: "Win Rate", value: "—" },
    { label: "Best Hand", value: "—" },
  ] as const;

  return (
    <Card className="border-border">
      <CardHeader>
        <CardTitle className="text-sm">Statistics</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4">
          {stats.map((stat) => (
            <div key={stat.label} className="space-y-1">
              <p className="text-xs text-muted-foreground">{stat.label}</p>
              {isLoading ? (
                <Skeleton className="h-5 w-16" />
              ) : (
                <p className="text-sm font-semibold font-mono">{stat.value}</p>
              )}
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
