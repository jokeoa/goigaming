"use client";

import { useParams } from "next/navigation";
import { RouletteGameView } from "@/components/roulette/roulette-game-view";
import { PageHeader } from "@/components/shared/page-header";
import { useRouletteTable } from "@/hooks/use-roulette";

export default function RouletteTablePage() {
  const params = useParams();
  const tableId = typeof params.id === "string" ? params.id : "";
  const { data: table } = useRouletteTable(tableId);

  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <PageHeader
        title={table?.name ?? "Roulette Table"}
        description={
          table
            ? `Min: $${table.min_bet} / Max: $${table.max_bet}`
            : "Loading..."
        }
      />
      {tableId && <RouletteGameView tableId={tableId} />}
    </div>
  );
}
