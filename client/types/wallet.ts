import { z } from "zod";

export const walletBalanceSchema = z.object({
  user_id: z.string(),
  balance: z.string(),
  updated_at: z.string(),
});

export type WalletBalance = z.infer<typeof walletBalanceSchema>;

export const transactionSchema = z.object({
  id: z.string(),
  wallet_id: z.string(),
  amount: z.string(),
  balance_after: z.string(),
  reference_type: z.string(),
  reference_id: z.string().nullable(),
  created_at: z.string(),
});

export type Transaction = z.infer<typeof transactionSchema>;

export const amountInputSchema = z.object({
  amount: z.string().refine(
    (val) => {
      const num = Number.parseFloat(val);
      return !Number.isNaN(num) && num > 0;
    },
    { message: "Amount must be a positive number" },
  ),
});

export type AmountInput = z.infer<typeof amountInputSchema>;
