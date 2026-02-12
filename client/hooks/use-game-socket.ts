"use client";

import { useCallback, useEffect, useRef } from "react";
import { WS_URL } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";
import { useGameStore } from "@/stores/game-store";
import type { WSMessage } from "@/types/ws";

export function useGameSocket(gameId: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const token = useAuthStore((s) => s.token);
  const setConnected = useGameStore((s) => s.setConnected);
  const reset = useGameStore((s) => s.reset);

  const connect = useCallback(() => {
    if (!gameId || !token) return;

    const url = `${WS_URL}?token=${encodeURIComponent(token)}&game_id=${encodeURIComponent(gameId)}`;
    const ws = new WebSocket(url);

    ws.onopen = () => {
      setConnected(true);
    };

    ws.onclose = () => {
      setConnected(false);
    };

    ws.onerror = () => {
      setConnected(false);
    };

    ws.onmessage = (event) => {
      try {
        const _message: WSMessage = JSON.parse(event.data);
      } catch {
        // ignore malformed messages
      }
    };

    wsRef.current = ws;
  }, [gameId, token, setConnected]);

  const disconnect = useCallback(() => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
    reset();
  }, [reset]);

  useEffect(() => {
    connect();
    return () => {
      disconnect();
    };
  }, [connect, disconnect]);

  return { disconnect };
}
