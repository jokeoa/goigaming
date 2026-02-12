"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getCurrentRound,
  getMyBets,
  getRouletteTable,
  getRound,
  getRoundHistory,
  listRouletteTables,
  placeBet,
} from "@/lib/api/roulette";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function useRouletteTables() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteTables,
    queryFn: listRouletteTables,
    enabled: isAuthenticated,
  });
}

export function useRouletteTable(id: string) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteTable(id),
    queryFn: () => getRouletteTable(id),
    enabled: isAuthenticated && !!id,
  });
}

export function useCurrentRound(tableId: string) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteCurrentRound(tableId),
    queryFn: () => getCurrentRound(tableId),
    enabled: isAuthenticated && !!tableId,
    refetchInterval: 5000,
  });
}

export function usePlaceBet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      tableId,
      body,
    }: {
      tableId: string;
      body: {
        round_id: string;
        bet_type: string;
        bet_value: string;
        amount: string;
      };
    }) => placeBet(tableId, body),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.rouletteCurrentRound(variables.tableId),
      });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.walletBalance,
      });
    },
  });
}

export function useRoundHistory(
  tableId: string,
  params?: { limit?: number; offset?: number },
) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteHistory(tableId, params),
    queryFn: () => getRoundHistory(tableId, params),
    enabled: isAuthenticated && !!tableId,
  });
}

export function useMyRouletteBets(params?: {
  limit?: number;
  offset?: number;
}) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteMyBets(params),
    queryFn: () => getMyBets(params),
    enabled: isAuthenticated,
  });
}

export function useRouletteRound(roundId: string) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.rouletteRound(roundId),
    queryFn: () => getRound(roundId),
    enabled: isAuthenticated && !!roundId,
  });
}
