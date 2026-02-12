"use client";

import { ActionBar } from "@/components/game/action-bar";
import { CommunityCards } from "@/components/game/community-cards";
import { PlayingCard } from "@/components/game/playing-card";
import { PotDisplay } from "@/components/game/pot-display";
import { Seat } from "@/components/game/seat";
import { TableInfo } from "@/components/game/table-info";
import { useGameStore } from "@/stores/game-store";

export function PokerTableView() {
  const tableState = useGameStore((s) => s.tableState);
  const holeCards = useGameStore((s) => s.holeCards);
  const isConnected = useGameStore((s) => s.isConnected);

  if (!tableState) {
    return (
      <div className="flex items-center justify-center rounded-lg border border-dashed border-border bg-card/50 p-12 text-sm text-muted-foreground">
        Connecting to table...
      </div>
    );
  }

  const maxPlayers = 6;
  const players = tableState.players ?? [];
  const communityCards = tableState.community_cards ?? [];
  const seats = Array.from({ length: maxPlayers }).map((_, i) =>
    players.find((p) => p.seat_number === i),
  );

  return (
    <div className="space-y-4">
      <TableInfo
        name={tableState.name}
        smallBlind={tableState.small_blind}
        bigBlind={tableState.big_blind}
        maxPlayers={maxPlayers}
        playerCount={players.length}
        stage={tableState.stage}
        isConnected={isConnected}
      />

      <div className="relative rounded-2xl border border-border bg-felt/20 p-6">
        <div className="grid grid-cols-3 gap-4">
          {/* Top row seats */}
          <div className="flex justify-center">
            <Seat
              player={seats[1]}
              position={1}
              isCurrentTurn={seats[1]?.user_id === tableState.current_turn}
            />
          </div>
          <div className="flex justify-center">
            <Seat
              player={seats[2]}
              position={2}
              isCurrentTurn={seats[2]?.user_id === tableState.current_turn}
            />
          </div>
          <div className="flex justify-center">
            <Seat
              player={seats[3]}
              position={3}
              isCurrentTurn={seats[3]?.user_id === tableState.current_turn}
            />
          </div>

          {/* Middle: left seat, community cards, right seat */}
          <div className="flex items-center justify-center">
            <Seat
              player={seats[0]}
              position={0}
              isCurrentTurn={seats[0]?.user_id === tableState.current_turn}
            />
          </div>
          <div className="flex flex-col items-center justify-center gap-3">
            <CommunityCards cards={communityCards} />
            <PotDisplay amount={tableState.pot} />
          </div>
          <div className="flex items-center justify-center">
            {maxPlayers > 4 && (
              <Seat
                player={seats[4]}
                position={4}
                isCurrentTurn={seats[4]?.user_id === tableState.current_turn}
              />
            )}
          </div>

          {/* Bottom row */}
          <div />
          <div className="flex justify-center">
            {maxPlayers > 5 && (
              <Seat
                player={seats[5]}
                position={5}
                isCurrentTurn={seats[5]?.user_id === tableState.current_turn}
              />
            )}
          </div>
          <div />
        </div>
      </div>

      {/* Hole cards */}
      {holeCards.length > 0 && (
        <div className="flex items-center justify-center gap-2">
          <span className="text-xs text-muted-foreground mr-2">Your hand:</span>
          {holeCards.map((card, i) => (
            <PlayingCard
              key={`hole-${card.suit}-${card.rank}-${i}`}
              card={card}
            />
          ))}
        </div>
      )}

      {/* Action bar */}
      <ActionBar isMyTurn={false} minBet={tableState.big_blind} />
    </div>
  );
}
