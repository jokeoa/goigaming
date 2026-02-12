import { PlayingCard } from "@/components/game/playing-card";
import type { Card } from "@/types/game";

type CommunityCardsProps = {
  readonly cards: readonly Card[];
};

export function CommunityCards({ cards }: CommunityCardsProps) {
  const safeCards = cards ?? [];
  const placeholders = 5 - safeCards.length;

  return (
    <div className="flex items-center justify-center gap-2">
      {safeCards.map((card, i) => (
        <PlayingCard key={`${card.suit}-${card.rank}-${i}`} card={card} />
      ))}
      {Array.from({ length: placeholders }).map((_, i) => (
        <div
          key={`placeholder-${i}`}
          className="flex h-16 w-11 items-center justify-center rounded-md border border-dashed border-border bg-card/30"
        />
      ))}
    </div>
  );
}
