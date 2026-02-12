"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { useRouletteStore } from "@/stores/roulette-store";
import type { RouletteBetType } from "@/types/roulette";

const amountPresets = ["1", "5", "10", "25", "50", "100"] as const;

const insideBets: { type: RouletteBetType; label: string }[] = [
  { type: "straight", label: "Straight" },
  { type: "split", label: "Split" },
  { type: "street", label: "Street" },
  { type: "corner", label: "Corner" },
  { type: "line", label: "Line" },
];

const outsideBets: { type: RouletteBetType; label: string }[] = [
  { type: "red", label: "Red" },
  { type: "black", label: "Black" },
  { type: "odd", label: "Odd" },
  { type: "even", label: "Even" },
  { type: "high", label: "High (19-36)" },
  { type: "low", label: "Low (1-18)" },
  { type: "dozen", label: "Dozen" },
  { type: "column", label: "Column" },
];

export function BetSelector() {
  const [selectedType, setSelectedType] = useState<RouletteBetType | null>(
    null,
  );
  const [betValue, setBetValue] = useState("");
  const selectedBetAmount = useRouletteStore((s) => s.selectedBetAmount);
  const setSelectedBetAmount = useRouletteStore((s) => s.setSelectedBetAmount);
  const addPendingBet = useRouletteStore((s) => s.addPendingBet);

  const needsValue =
    selectedType === "straight" ||
    selectedType === "split" ||
    selectedType === "street" ||
    selectedType === "corner" ||
    selectedType === "line" ||
    selectedType === "dozen" ||
    selectedType === "column";

  const handleAddBet = () => {
    if (!selectedType) return;

    const value = needsValue ? betValue : selectedType;
    if (needsValue && !betValue) return;

    addPendingBet({
      bet_type: selectedType,
      bet_value: value,
      amount: selectedBetAmount,
    });

    setSelectedType(null);
    setBetValue("");
  };

  return (
    <div className="space-y-4 rounded-lg border border-border bg-card p-4">
      <div className="space-y-2">
        <h3 className="text-sm font-medium">Bet Amount</h3>
        <div className="flex flex-wrap gap-1.5">
          {amountPresets.map((amount) => (
            <Button
              key={amount}
              size="sm"
              variant={selectedBetAmount === amount ? "default" : "outline"}
              className="h-7 px-2.5 text-xs font-mono"
              onClick={() => setSelectedBetAmount(amount)}
            >
              ${amount}
            </Button>
          ))}
          <Input
            type="number"
            placeholder="Custom"
            value={
              amountPresets.includes(
                selectedBetAmount as (typeof amountPresets)[number],
              )
                ? ""
                : selectedBetAmount
            }
            onChange={(e) => setSelectedBetAmount(e.target.value)}
            className="h-7 w-20 text-xs font-mono"
          />
        </div>
      </div>

      <div className="space-y-2">
        <h3 className="text-sm font-medium">Inside Bets</h3>
        <div className="flex flex-wrap gap-1.5">
          {insideBets.map((bet) => (
            <Button
              key={bet.type}
              size="sm"
              variant={selectedType === bet.type ? "default" : "outline"}
              className="h-7 text-xs"
              onClick={() => setSelectedType(bet.type)}
            >
              {bet.label}
            </Button>
          ))}
        </div>
      </div>

      <div className="space-y-2">
        <h3 className="text-sm font-medium">Outside Bets</h3>
        <div className="flex flex-wrap gap-1.5">
          {outsideBets.map((bet) => (
            <Button
              key={bet.type}
              size="sm"
              variant={selectedType === bet.type ? "default" : "outline"}
              className={cn(
                "h-7 text-xs",
                bet.type === "red" &&
                  selectedType === "red" &&
                  "bg-red-600 hover:bg-red-700",
                bet.type === "black" &&
                  selectedType === "black" &&
                  "bg-zinc-800 hover:bg-zinc-900",
              )}
              onClick={() => setSelectedType(bet.type)}
            >
              {bet.label}
            </Button>
          ))}
        </div>
      </div>

      {selectedType && needsValue && (
        <div className="space-y-2">
          <h3 className="text-sm font-medium">
            Bet Value
            {selectedType === "straight" && " (0-36)"}
            {selectedType === "dozen" && " (1st, 2nd, 3rd)"}
            {selectedType === "column" && " (1st, 2nd, 3rd)"}
          </h3>
          <Input
            type="text"
            placeholder="Enter value"
            value={betValue}
            onChange={(e) => setBetValue(e.target.value)}
            className="h-8 w-32 text-sm"
          />
        </div>
      )}

      {selectedType && (
        <Button
          size="sm"
          onClick={handleAddBet}
          disabled={needsValue && !betValue}
        >
          Add {selectedType} bet (${selectedBetAmount})
        </Button>
      )}
    </div>
  );
}
