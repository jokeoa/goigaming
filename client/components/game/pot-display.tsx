import { formatCurrency } from "@/lib/utils";

type PotDisplayProps = {
  readonly amount: string;
};

export function PotDisplay({ amount }: PotDisplayProps) {
  return (
    <div className="flex items-center justify-center gap-2 rounded-full border border-gold/30 bg-gold/10 px-4 py-1.5">
      <span className="text-xs text-muted-foreground">Pot</span>
      <span className="text-sm font-bold font-mono text-gold">
        {formatCurrency(amount)}
      </span>
    </div>
  );
}
