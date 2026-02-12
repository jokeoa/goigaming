"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type ActionBarProps = {
  readonly isMyTurn: boolean;
  readonly minBet?: number;
  readonly onFold?: () => void;
  readonly onCheck?: () => void;
  readonly onCall?: () => void;
  readonly onRaise?: (amount: number) => void;
};

export function ActionBar({
  isMyTurn,
  minBet = 0,
  onFold,
  onCheck,
  onCall,
  onRaise,
}: ActionBarProps) {
  const [raiseAmount, setRaiseAmount] = useState("");

  const handleRaise = () => {
    const amount = Number.parseFloat(raiseAmount);
    if (!Number.isNaN(amount) && amount > 0) {
      onRaise?.(amount);
      setRaiseAmount("");
    }
  };

  return (
    <div className="flex flex-wrap items-center justify-center gap-2 rounded-lg border border-border bg-card p-3">
      <Button
        variant="destructive"
        size="sm"
        disabled={!isMyTurn}
        onClick={onFold}
      >
        Fold
      </Button>
      <Button
        variant="outline"
        size="sm"
        disabled={!isMyTurn}
        onClick={onCheck}
      >
        Check
      </Button>
      <Button size="sm" disabled={!isMyTurn} onClick={onCall}>
        Call
      </Button>
      <div className="flex items-center gap-1">
        <Input
          type="number"
          placeholder={String(minBet)}
          value={raiseAmount}
          onChange={(e) => setRaiseAmount(e.target.value)}
          className="h-8 w-20 font-mono text-sm"
          disabled={!isMyTurn}
        />
        <Button
          size="sm"
          variant="outline"
          disabled={!isMyTurn}
          onClick={handleRaise}
        >
          Raise
        </Button>
      </div>
    </div>
  );
}
