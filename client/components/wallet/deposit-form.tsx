"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowDownToLine } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDeposit } from "@/hooks/use-wallet";
import { type AmountInput, amountInputSchema } from "@/types/wallet";

export function DepositForm() {
  const deposit = useDeposit();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AmountInput>({
    resolver: zodResolver(amountInputSchema),
  });

  const onSubmit = (data: AmountInput) => {
    deposit.mutate(data.amount, {
      onSuccess: () => {
        toast.success("Deposit successful");
        reset();
      },
      onError: (error) => {
        toast.error(error instanceof Error ? error.message : "Deposit failed");
      },
    });
  };

  return (
    <Card className="border-border">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-sm">
          <ArrowDownToLine className="h-4 w-4 text-primary" />
          Deposit
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
          <div className="space-y-2">
            <Label htmlFor="deposit-amount">Amount (USD)</Label>
            <Input
              id="deposit-amount"
              placeholder="100.00"
              className="font-mono"
              {...register("amount")}
            />
            {errors.amount && (
              <p className="text-xs text-destructive">
                {errors.amount.message}
              </p>
            )}
          </div>
          <Button type="submit" className="w-full" disabled={deposit.isPending}>
            {deposit.isPending ? "Processing..." : "Deposit"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
