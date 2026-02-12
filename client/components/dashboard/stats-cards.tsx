"use client";

import { CreditCard, Trophy, Users, Wallet } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useCurrentUser } from "@/hooks/use-user";
import { useWalletBalance } from "@/hooks/use-wallet";
import { formatCurrency } from "@/lib/utils";

export function StatsCards() {
  const { data: wallet, isLoading: walletLoading } = useWalletBalance();
  const { data: user, isLoading: userLoading } = useCurrentUser();

  const stats = [
    {
      label: "Balance",
      value: wallet ? formatCurrency(wallet.balance) : "$0.00",
      icon: Wallet,
      loading: walletLoading,
    },
    {
      label: "Player",
      value: user?.username ?? "—",
      icon: Users,
      loading: userLoading,
    },
    {
      label: "Games Played",
      value: "0",
      icon: Trophy,
      loading: false,
    },
    {
      label: "Account",
      value: "Active",
      icon: CreditCard,
      loading: false,
    },
  ] as const;

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat) => (
        <Card key={stat.label} className="border-border">
          <CardHeader className="flex flex-row items-center justify-between pb-2">
            <CardTitle className="text-xs font-medium text-muted-foreground">
              {stat.label}
            </CardTitle>
            <stat.icon className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            {stat.loading ? (
              <Skeleton className="h-6 w-24" />
            ) : (
              <p className="text-lg font-semibold font-mono">{stat.value}</p>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
