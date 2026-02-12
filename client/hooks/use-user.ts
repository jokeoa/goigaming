"use client";

import { useQuery } from "@tanstack/react-query";
import { getMe } from "@/lib/api/user";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function useCurrentUser() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.user,
    queryFn: getMe,
    enabled: isAuthenticated,
  });
}
