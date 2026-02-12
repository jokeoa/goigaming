"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getPokerTable,
  getPokerTableState,
  joinPokerTable,
  leavePokerTable,
} from "@/lib/api/game";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function usePokerTable(id: string) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.pokerTable(id),
    queryFn: () => getPokerTable(id),
    enabled: isAuthenticated && !!id,
  });
}

export function usePokerTableState(id: string) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.pokerTableState(id),
    queryFn: () => getPokerTableState(id),
    enabled: isAuthenticated && !!id,
  });
}

export function useJoinPokerTable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      tableId,
      seat_number,
      buy_in,
    }: {
      tableId: string;
      seat_number: number;
      buy_in: string;
    }) => joinPokerTable(tableId, { seat_number, buy_in }),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.pokerTableState(variables.tableId),
      });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.walletBalance,
      });
    },
  });
}

export function useLeavePokerTable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ tableId }: { tableId: string }) => leavePokerTable(tableId),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.pokerTableState(variables.tableId),
      });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.walletBalance,
      });
    },
  });
}
