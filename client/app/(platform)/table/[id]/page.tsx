"use client";

import { useParams } from "next/navigation";
import { PokerTableView } from "@/components/game/poker-table";
import { PageHeader } from "@/components/shared/page-header";
import { useGameSocket } from "@/hooks/use-game-socket";

export default function TablePage() {
  const params = useParams();
  const tableId = typeof params.id === "string" ? params.id : null;
  useGameSocket(tableId);

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <PageHeader
        title="Poker Table"
        description={tableId ? `Table: ${tableId}` : "Loading..."}
      />
      <PokerTableView />
    </div>
  );
}
