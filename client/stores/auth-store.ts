import { create } from "zustand";
import { persist } from "zustand/middleware";

type AuthState = {
  readonly token: string | null;
  readonly isAuthenticated: boolean;
};

type AuthActions = {
  readonly setToken: (token: string) => void;
  readonly logout: () => void;
};

export const useAuthStore = create<AuthState & AuthActions>()(
  persist(
    (set) => ({
      token: null,
      isAuthenticated: false,
      setToken: (token: string) => set({ token, isAuthenticated: true }),
      logout: () => set({ token: null, isAuthenticated: false }),
    }),
    {
      name: "goi-auth-storage",
    },
  ),
);
