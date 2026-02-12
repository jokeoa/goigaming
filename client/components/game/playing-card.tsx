import { cn } from "@/lib/utils";
import type { Card, Suit } from "@/types/game";

type PlayingCardProps = {
  readonly card?: Card;
  readonly faceDown?: boolean;
  readonly className?: string;
};

const suitSymbols: Record<Suit, string> = {
  h: "\u2665",
  d: "\u2666",
  c: "\u2663",
  s: "\u2660",
} as const;

const suitColors: Record<Suit, string> = {
  h: "text-suit-hearts",
  d: "text-suit-diamonds",
  c: "text-suit-clubs",
  s: "text-suit-spades",
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
