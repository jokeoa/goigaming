import { create } from "zustand";
import type { RouletteBetType } from "@/types/roulette";

type PendingBet = {
  readonly id: string;
  readonly bet_type: RouletteBetType;
  readonly bet_value: string;
  readonly amount: string;
};

type RouletteState = {
  readonly pendingBets: readonly PendingBet[];
  readonly selectedBetAmount: string;
};

type RouletteActions = {
  readonly addPendingBet: (bet: Omit<PendingBet, "id">) => void;
  readonly removePendingBet: (id: string) => void;
  readonly clearPendingBets: () => void;
  readonly setSelectedBetAmount: (amount: string) => void;
};

const initialState: RouletteState = {
  pendingBets: [],
  selectedBetAmount: "5",
};

export const useRouletteStore = create<RouletteState & RouletteActions>()(
  (set) => ({
    ...initialState,
    addPendingBet: (bet) =>
      set((state) => ({
        pendingBets: [
          ...state.pendingBets,
          { ...bet, id: crypto.randomUUID() },
        ],
      })),
    removePendingBet: (id) =>
      set((state) => ({
        pendingBets: state.pendingBets.filter((b) => b.id !== id),
      })),
    clearPendingBets: () => set({ pendingBets: [] }),
    setSelectedBetAmount: (selectedBetAmount) => set({ selectedBetAmount }),
  }),
);
