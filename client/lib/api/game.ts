import type { PokerTable } from "@/types/game";
import { api } from "./client";

export function listTables(): Promise<PokerTable[]> {
  return api.get<PokerTable[]>("/tables");
}

export function getTable(id: string): Promise<PokerTable> {
  return api.get<PokerTable>(`/tables/${id}`);
}
