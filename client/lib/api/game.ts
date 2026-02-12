import type { PokerPlayer, PokerTable, WSTableState } from "@/types/game";
import { api } from "./client";

export function listPokerTables(): Promise<PokerTable[]> {
  return api.get<PokerTable[]>("/poker/tables");
}

export function getPokerTable(id: string): Promise<PokerTable> {
  return api.get<PokerTable>(`/poker/tables/${id}`);
}

export function joinPokerTable(
  tableId: string,
  body: { seat_number: number; buy_in: string },
): Promise<PokerPlayer> {
  return api.post<PokerPlayer>(`/poker/tables/${tableId}/join`, body);
}

export function leavePokerTable(tableId: string): Promise<{ message: string }> {
  return api.post<{ message: string }>(`/poker/tables/${tableId}/leave`);
}

export function getPokerTableState(tableId: string): Promise<WSTableState> {
  return api.get<WSTableState>(`/poker/tables/${tableId}/state`);
}
