"use client";

import { ActionBar } from "@/components/game/action-bar";
import { CommunityCards } from "@/components/game/community-cards";
import { PlayingCard } from "@/components/game/playing-card";
import { PotDisplay } from "@/components/game/pot-display";
import { Seat } from "@/components/game/seat";
import { TableInfo } from "@/components/game/table-info";
import { useGameStore } from "@/stores/game-store";
import type {
  Card,
  GameStage,
  PlayerInfo,
  PokerTable as PokerTableType,
} from "@/types/game";

const mockTable: PokerTableType = {
  id: "mock-1",
  name: "Table Alpha",
  maxPlayers: 6,
  currentPlayers: 3,
  smallBlind: 5,
  bigBlind: 10,
  minBuyIn: 100,
  maxBuyIn: 1000,
  stage: "flop",
};

const mockPlayers: readonly PlayerInfo[] = [
  {
    id: "p1",
    username: "Alice",
    seatIndex: 0,
    chips: 450,
    currentBet: 20,
    isFolded: false,
    isDealer: true,
    isActive: true,
  },
  {
    id: "p2",
    username: "Bob",
    seatIndex: 2,
    chips: 780,
    currentBet: 20,
    isFolded: false,
    isDealer: false,
    isActive: true,
  },
  {
    id: "p3",
    username: "Charlie",
    seatIndex: 4,
    chips: 200,
    currentBet: 0,
    isFolded: true,
    isDealer: false,
    isActive: false,
  },
] as const;

const mockCommunityCards: readonly Card[] = [
  { suit: "hearts", rank: "A" },
  { suit: "spades", rank: "K" },
  { suit: "diamonds", rank: "7" },
] as const;

const mockHoleCards: readonly Card[] = [
  { suit: "clubs", rank: "Q" },
  { suit: "hearts", rank: "J" },
] as const;

export function PokerTableView() {
  const tableState = useGameStore((s) => s.tableState);
  const holeCards = useGameStore((s) => s.holeCards);
  const isConnected = useGameStore((s) => s.isConnected);

  const table = tableState?.table ?? mockTable;
  const players = tableState?.players ?? mockPlayers;
  const communityCards = tableState?.communityCards ?? mockCommunityCards;
  const pot = tableState?.pot ?? 60;
  const stage: GameStage = tableState?.stage ?? "flop";
  const currentTurn = tableState?.currentTurn ?? "p1";
  const displayHoleCards = holeCards.length > 0 ? holeCards : mockHoleCards;

  const seats = Array.from({ length: table.maxPlayers }).map((_, i) =>
    players.find((p) => p.seatIndex === i),
  );

  return (
    <div className="space-y-4">
      <TableInfo table={table} stage={stage} isConnected={isConnected} />

      <div className="relative rounded-2xl border border-border bg-felt/20 p-6">
        <div className="grid grid-cols-3 gap-4">
          {/* Top row seats */}
          <div className="flex justify-center">
            <Seat
              player={seats[1]}
              position={1}
              isCurrentTurn={seats[1]?.id === currentTurn}
            />
          </div>
          <div className="flex justify-center">
            <Seat
              player={seats[2]}
              position={2}
              isCurrentTurn={seats[2]?.id === currentTurn}
            />
          </div>
          <div className="flex justify-center">
            <Seat
              player={seats[3]}
              position={3}
              isCurrentTurn={seats[3]?.id === currentTurn}
            />
          </div>

          {/* Middle: left seat, community cards, right seat */}
          <div className="flex items-center justify-center">
            <Seat
              player={seats[0]}
              position={0}
              isCurrentTurn={seats[0]?.id === currentTurn}
            />
          </div>
          <div className="flex flex-col items-center justify-center gap-3">
            <CommunityCards cards={communityCards} />
            <PotDisplay amount={pot} />
          </div>
          <div className="flex items-center justify-center">
            {table.maxPlayers > 4 && (
              <Seat
                player={seats[4]}
                position={4}
                isCurrentTurn={seats[4]?.id === currentTurn}
              />
            )}
          </div>

          {/* Bottom row */}
          <div />
          <div className="flex justify-center">
            {table.maxPlayers > 5 && (
              <Seat
                player={seats[5]}
                position={5}
                isCurrentTurn={seats[5]?.id === currentTurn}
              />
            )}
          </div>
          <div />
        </div>
      </div>

      {/* Hole cards */}
      <div className="flex items-center justify-center gap-2">
        <span className="text-xs text-muted-foreground mr-2">Your hand:</span>
        {displayHoleCards.map((card, i) => (
          <PlayingCard
            key={`hole-${card.suit}-${card.rank}-${i}`}
            card={card}
          />
        ))}
      </div>

      {/* Action bar */}
      <ActionBar isMyTurn={currentTurn === "p1"} minBet={table.bigBlind} />
    </div>
  );
}
