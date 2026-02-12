import { cn } from "@/lib/utils";
import type { Card } from "@/types/game";

type PlayingCardProps = {
  readonly card?: Card;
  readonly faceDown?: boolean;
  readonly className?: string;
};

const suitSymbols = {
  hearts: "\u2665",
  diamonds: "\u2666",
  clubs: "\u2663",
  spades: "\u2660",
} as const;

const suitColors = {
  hearts: "text-suit-hearts",
  diamonds: "text-suit-diamonds",
  clubs: "text-suit-clubs",
  spades: "text-suit-spades",
} as const;

export function PlayingCard({ card, faceDown, className }: PlayingCardProps) {
  if (faceDown || !card) {
    return (
      <div
        className={cn(
          "flex h-16 w-11 items-center justify-center rounded-md border border-border bg-felt text-xs font-bold text-primary",
          className,
        )}
      >
        ?
      </div>
    );
  }

  return (
    <div
      className={cn(
        "flex h-16 w-11 flex-col items-center justify-center rounded-md border border-border bg-card text-xs font-bold",
        suitColors[card.suit],
        className,
      )}
    >
      <span className="text-sm leading-none">{card.rank}</span>
      <span className="text-base leading-none">{suitSymbols[card.suit]}</span>
    </div>
  );
}
