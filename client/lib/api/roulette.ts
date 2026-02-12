import type {
  RouletteBet,
  RouletteRound,
  RouletteTable,
} from "@/types/roulette";
import { api } from "./client";

export function listRouletteTables(): Promise<RouletteTable[]> {
  return api.get<RouletteTable[]>("/roulette/tables");
}

export function getRouletteTable(id: string): Promise<RouletteTable> {
  return api.get<RouletteTable>(`/roulette/tables/${id}`);
}

export function getCurrentRound(tableId: string): Promise<RouletteRound> {
  return api.get<RouletteRound>(`/roulette/tables/${tableId}/current-round`);
}

export function placeBet(
  tableId: string,
  body: {
    round_id: string;
    bet_type: string;
    bet_value: string;
    amount: string;
  },
): Promise<RouletteBet> {
  return api.post<RouletteBet>(`/roulette/tables/${tableId}/bets`, body);
}

export function getRoundHistory(
  tableId: string,
  params?: { limit?: number; offset?: number },
): Promise<RouletteRound[]> {
  return api.get<RouletteRound[]>(
    `/roulette/tables/${tableId}/history`,
    params,
  );
}

export function getMyBets(params?: {
  limit?: number;
  offset?: number;
}): Promise<RouletteBet[]> {
  return api.get<RouletteBet[]>("/roulette/bets/me", params);
}

export function getRound(roundId: string): Promise<RouletteRound> {
  return api.get<RouletteRound>(`/roulette/rounds/${roundId}`);
}
