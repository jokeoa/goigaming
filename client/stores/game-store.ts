import { create } from "zustand";
import type { Card, WSTableState } from "@/types/game";

type GameState = {
  readonly tableState: WSTableState | null;
  readonly holeCards: readonly Card[];
  readonly isConnected: boolean;
};

type GameActions = {
  readonly setTableState: (state: WSTableState) => void;
  readonly setHoleCards: (cards: readonly Card[]) => void;
  readonly setConnected: (connected: boolean) => void;
  readonly reset: () => void;
};

const initialState: GameState = {
  tableState: null,
  holeCards: [],
  isConnected: false,
};

export const useGameStore = create<GameState & GameActions>()((set) => ({
  ...initialState,
  setTableState: (tableState) => set({ tableState }),
  setHoleCards: (holeCards) => set({ holeCards }),
  setConnected: (isConnected) => set({ isConnected }),
  reset: () => set(initialState),
}));
