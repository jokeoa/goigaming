"use client";

import { useQuery } from "@tanstack/react-query";
import { listPokerTables } from "@/lib/api/game";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function usePokerTables() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.pokerTables,
    queryFn: listPokerTables,
    enabled: isAuthenticated,
  });
}
