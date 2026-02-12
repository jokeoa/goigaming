"use client";

import { Gamepad2 } from "lucide-react";
import { EmptyState } from "@/components/shared/empty-state";
import { Skeleton } from "@/components/ui/skeleton";
import { usePokerTables } from "@/hooks/use-tables";
import { TableCard } from "./table-card";

export function TableList() {
  const { data: tables, isLoading, error } = usePokerTables();

  if (isLoading) {
    return (
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={`table-skeleton-${i}`} className="h-48" />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <EmptyState
        icon={<Gamepad2 className="h-8 w-8" />}
        title="Failed to load tables"
        description="Tables will appear here once the game server is running."
      />
    );
  }

  if (!tables || tables.length === 0) {
    return (
      <EmptyState
        icon={<Gamepad2 className="h-8 w-8" />}
        title="No tables available"
        description="Check back soon for open tables."
      />
    );
  }

  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {tables.map((table) => (
        <TableCard key={table.id} table={table} />
      ))}
    </div>
  );
}
