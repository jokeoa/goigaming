import type { Transaction, WalletBalance } from "@/types/wallet";
import { api } from "./client";

export function getBalance(): Promise<WalletBalance> {
  return api.get<WalletBalance>("/wallet/balance");
}

export function deposit(amount: string): Promise<WalletBalance> {
  return api.post<WalletBalance>("/wallet/deposit", { amount });
}

export function withdraw(amount: string): Promise<WalletBalance> {
  return api.post<WalletBalance>("/wallet/withdraw", { amount });
}

export function getTransactions(params?: {
  readonly limit?: number;
  readonly offset?: number;
}): Promise<Transaction[]> {
  return api.get<Transaction[]>(
    "/wallet/transactions",
    params as Record<string, number | undefined>,
  );
}
