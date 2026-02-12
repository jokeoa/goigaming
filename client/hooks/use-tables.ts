"use client";

import { useQuery } from "@tanstack/react-query";
import { listTables } from "@/lib/api/game";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function useTables() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.tables,
    queryFn: listTables,
    enabled: isAuthenticated,
  });
}
