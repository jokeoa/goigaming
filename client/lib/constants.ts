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
  tables: ["tables"] as const,
  table: (id: string) => ["tables", id] as const,
} as const;
