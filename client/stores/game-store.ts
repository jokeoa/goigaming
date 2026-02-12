import { create } from "zustand";
import type { Card, GameStage, PlayerInfo, TableState } from "@/types/game";

type GameState = {
  readonly tableState: TableState | null;
  readonly holeCards: readonly Card[];
  readonly isConnected: boolean;
};

type GameActions = {
  readonly setTableState: (state: TableState) => void;
  readonly setHoleCards: (cards: readonly Card[]) => void;
  readonly setConnected: (connected: boolean) => void;
  readonly updatePlayers: (players: readonly PlayerInfo[]) => void;
  readonly updateCommunityCards: (cards: readonly Card[]) => void;
  readonly updatePot: (pot: number) => void;
  readonly updateStage: (stage: GameStage) => void;
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
  updatePlayers: (players) =>
    set((state) => ({
      tableState: state.tableState ? { ...state.tableState, players } : null,
    })),
  updateCommunityCards: (communityCards) =>
    set((state) => ({
      tableState: state.tableState
        ? { ...state.tableState, communityCards }
        : null,
    })),
  updatePot: (pot) =>
    set((state) => ({
      tableState: state.tableState ? { ...state.tableState, pot } : null,
    })),
  updateStage: (stage) =>
    set((state) => ({
      tableState: state.tableState ? { ...state.tableState, stage } : null,
    })),
  reset: () => set(initialState),
}));
