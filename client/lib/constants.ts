export const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

export const TOKEN_STORAGE_KEY = "goi-auth-token";

export const QUERY_KEYS = {
  user: ["user"] as const,
  walletBalance: ["wallet", "balance"] as const,
  transactions: (params?: { limit?: number; offset?: number }) =>
    ["wallet", "transactions", params] as const,
  pokerTables: ["poker", "tables"] as const,
  pokerTable: (id: string) => ["poker", "tables", id] as const,
  pokerTableState: (id: string) => ["poker", "tables", id, "state"] as const,
  rouletteTables: ["roulette", "tables"] as const,
  rouletteTable: (id: string) => ["roulette", "tables", id] as const,
  rouletteCurrentRound: (tableId: string) =>
    ["roulette", "tables", tableId, "current-round"] as const,
  rouletteHistory: (
    tableId: string,
    params?: { limit?: number; offset?: number },
  ) => ["roulette", "tables", tableId, "history", params] as const,
  rouletteMyBets: (params?: { limit?: number; offset?: number }) =>
    ["roulette", "bets", "me", params] as const,
  rouletteRound: (id: string) => ["roulette", "rounds", id] as const,
} as const;
