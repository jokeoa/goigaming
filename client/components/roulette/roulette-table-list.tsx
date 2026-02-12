"use client";

import { CircleDot } from "lucide-react";
import { EmptyState } from "@/components/shared/empty-state";
import { Skeleton } from "@/components/ui/skeleton";
import { useRouletteTables } from "@/hooks/use-roulette";
import { RouletteTableCard } from "./roulette-table-card";

export function RouletteTableList() {
  const { data: tables, isLoading, error } = useRouletteTables();

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={`roulette-skeleton-${i}`} className="h-40" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <EmptyState
        icon={<CircleDot className="h-8 w-8" />}
        title="Failed to load roulette tables"
        description="Tables will appear here once the game server is running."
      />
    );
  }

  if (!tables || tables.length === 0) {
    return (
      <EmptyState
        icon={<CircleDot className="h-8 w-8" />}
        title="No roulette tables available"
        description="Check back soon for open tables."
      />
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {tables.map((table) => (
        <RouletteTableCard key={table.id} table={table} />
      ))}
    </div>
  );
}
