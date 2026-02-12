"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { deposit, getBalance, withdraw } from "@/lib/api/wallet";
import { QUERY_KEYS } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";

export function useWalletBalance() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);

  return useQuery({
    queryKey: QUERY_KEYS.walletBalance,
    queryFn: getBalance,
    enabled: isAuthenticated,
  });
}

export function useDeposit() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (amount: string) => deposit(amount),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.walletBalance });
      queryClient.invalidateQueries({ queryKey: ["wallet", "transactions"] });
    },
  });
}

export function useWithdraw() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (amount: string) => withdraw(amount),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.walletBalance });
      queryClient.invalidateQueries({ queryKey: ["wallet", "transactions"] });
    },
  });
}
