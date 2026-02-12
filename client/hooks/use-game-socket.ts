"use client";

import { useCallback, useEffect, useRef } from "react";
import { WS_URL } from "@/lib/constants";
import { useAuthStore } from "@/stores/auth-store";
import { useGameStore } from "@/stores/game-store";
import type { WSCardsDealt, WSTableState } from "@/types/game";
import type { WSMessage } from "@/types/ws";

export function useGameSocket(gameId: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const token = useAuthStore((s) => s.token);
  const setConnected = useGameStore((s) => s.setConnected);
  const setTableState = useGameStore((s) => s.setTableState);
  const setHoleCards = useGameStore((s) => s.setHoleCards);
  const reset = useGameStore((s) => s.reset);

  const handleMessage = useCallback(
    (message: WSMessage) => {
      switch (message.type) {
        case "table_state": {
          const payload = message.payload as WSTableState;
          setTableState(payload);
          break;
        }
        case "cards_dealt": {
          const payload = message.payload as WSCardsDealt;
          setHoleCards(payload.hole_cards);
          break;
        }
        case "player_joined":
        case "player_left":
        case "player_acted":
        case "community_cards":
        case "pot_updated":
        case "turn_changed":
        case "new_hand":
        case "hand_result": {
          const payload = message.payload as WSTableState;
          setTableState(payload);
          break;
        }
      }
    },
    [setTableState, setHoleCards],
  );

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
        const message: WSMessage = JSON.parse(event.data);
        handleMessage(message);
      } catch {
        // ignore malformed messages
      }
    };

    wsRef.current = ws;
  }, [gameId, token, setConnected, handleMessage]);

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
