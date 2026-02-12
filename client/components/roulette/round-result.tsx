import { cn } from "@/lib/utils";
import type { RouletteRound } from "@/types/roulette";

type RoundResultProps = {
  readonly round: RouletteRound | undefined;
};

const colorClasses = {
  red: "bg-red-600 text-white",
  black: "bg-zinc-800 text-white",
  green: "bg-emerald-600 text-white",
} as const;

export function RoundResult({ round }: RoundResultProps) {
  if (!round || round.result === null || !round.result_color) {
    return null;
  }

  const colorClass =
    colorClasses[round.result_color as keyof typeof colorClasses] ??
    colorClasses.green;

  return (
    <div className="flex flex-col items-center gap-2 py-4">
      <span className="text-xs text-muted-foreground uppercase tracking-wider">
        Result
      </span>
      <div
        className={cn(
          "flex h-20 w-20 items-center justify-center rounded-full text-3xl font-bold shadow-lg",
          colorClass,
        )}
      >
        {round.result}
      </div>
      <span className="text-xs text-muted-foreground capitalize">
        {round.result_color}
      </span>
    </div>
  );
}
