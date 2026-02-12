import { z } from "zod";

export type RouletteTableStatus = "active" | "inactive" | "maintenance";

export type RouletteTable = {
  readonly id: string;
  readonly name: string;
  readonly min_bet: string;
  readonly max_bet: string;
  readonly status: RouletteTableStatus;
  readonly created_at: string;
};

export type RouletteRound = {
  readonly id: string;
  readonly table_id: string;
  readonly round_number: number;
  readonly result: number | null;
  readonly result_color: string | null;
  readonly seed_hash: string | null;
  readonly seed_revealed: string | null;
  readonly betting_ends_at: string | null;
  readonly created_at: string;
  readonly settled_at: string | null;
};

export type RouletteBetStatus = "pending" | "won" | "lost";

export type RouletteBet = {
  readonly id: string;
  readonly round_id: string;
  readonly user_id: string;
  readonly bet_type: RouletteBetType;
  readonly bet_value: string;
  readonly amount: string;
  readonly payout: string;
  readonly status: RouletteBetStatus;
  readonly created_at: string;
};

export type RouletteBetType =
  | "straight"
  | "split"
  | "street"
  | "corner"
  | "line"
  | "dozen"
  | "column"
  | "red"
  | "black"
  | "odd"
  | "even"
  | "high"
  | "low";

export const placeBetSchema = z.object({
  round_id: z.string().uuid(),
  bet_type: z.enum([
    "straight",
    "split",
    "street",
    "corner",
    "line",
    "dozen",
    "column",
    "red",
    "black",
    "odd",
    "even",
    "high",
    "low",
  ]),
  bet_value: z.string().min(1),
  amount: z.string().regex(/^\d+(\.\d{1,2})?$/, "Invalid amount"),
});

export type PlaceBetInput = z.infer<typeof placeBetSchema>;
