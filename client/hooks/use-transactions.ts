"use client";

import { useQuery } from "@tanstack/react-query";
import { getTransactions } from "@/lib/api/wallet";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function useTransactions(params?: {
  readonly limit?: number;
  readonly offset?: number;
}) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.transactions(params),
    queryFn: () => getTransactions(params),
    enabled: isAuthenticated,
  });
}
