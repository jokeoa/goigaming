"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowUpFromLine } from "lucide-react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useWithdraw } from "@/hooks/use-wallet";
import { type AmountInput, amountInputSchema } from "@/types/wallet";

export function WithdrawForm() {
  const withdraw = useWithdraw();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AmountInput>({
    resolver: zodResolver(amountInputSchema),
  });

  const onSubmit = (data: AmountInput) => {
    withdraw.mutate(data.amount, {
      onSuccess: () => {
        toast.success("Withdrawal successful");
        reset();
      },
      onError: (error) => {
        toast.error(
          error instanceof Error ? error.message : "Withdrawal failed",
        );
      },
    });
  };

  return (
    <Card className="border-border">
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-sm">
          <ArrowUpFromLine className="h-4 w-4 text-destructive" />
          Withdraw
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
          <div className="space-y-2">
            <Label htmlFor="withdraw-amount">Amount (USD)</Label>
            <Input
              id="withdraw-amount"
              placeholder="50.00"
              className="font-mono"
              {...register("amount")}
            />
            {errors.amount && (
              <p className="text-xs text-destructive">
                {errors.amount.message}
              </p>
            )}
          </div>
          <Button
            type="submit"
            variant="outline"
            className="w-full"
            disabled={withdraw.isPending}
          >
            {withdraw.isPending ? "Processing..." : "Withdraw"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
