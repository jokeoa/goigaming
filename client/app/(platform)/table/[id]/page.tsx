"use client";

import { useParams } from "next/navigation";
import { useEffect } from "react";
import { PokerTableView } from "@/components/game/poker-table";
import { PageHeader } from "@/components/shared/page-header";
import { useGameSocket } from "@/hooks/use-game-socket";
import { usePokerTableState } from "@/hooks/use-poker";
import { useGameStore } from "@/stores/game-store";

export default function TablePage() {
  const params = useParams();
  const tableId = typeof params.id === "string" ? params.id : null;
  useGameSocket(tableId);

  const { data: initialState } = usePokerTableState(tableId ?? "");
  const setTableState = useGameStore((s) => s.setTableState);
  const tableState = useGameStore((s) => s.tableState);

  useEffect(() => {
    if (initialState && !tableState) {
      setTableState(initialState);
    }
  }, [initialState, tableState, setTableState]);

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
