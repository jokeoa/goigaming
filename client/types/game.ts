export type Suit = "s" | "h" | "d" | "c";

export type Rank =
  | "2"
  | "3"
  | "4"
  | "5"
  | "6"
  | "7"
  | "8"
  | "9"
  | "T"
  | "J"
  | "Q"
  | "K"
  | "A";

export type Card = {
  readonly suit: Suit;
  readonly rank: Rank;
};

export type GameStage =
  | "waiting"
  | "preflop"
  | "flop"
  | "turn"
  | "river"
  | "showdown"
  | "complete";

export type PlayerStatus = "active" | "sitting_out" | "all_in" | "folded";

export type TableStatus = "waiting" | "active" | "closed";

export type ActionType = "fold" | "check" | "call" | "raise" | "all_in" | "bet";

export type PokerTable = {
  readonly id: string;
  readonly name: string;
  readonly small_blind: string;
  readonly big_blind: string;
  readonly min_buy_in: string;
  readonly max_buy_in: string;
  readonly max_players: number;
  readonly status: TableStatus;
  readonly created_at: string;
};

export type PokerPlayer = {
  readonly id: string;
  readonly table_id: string;
  readonly user_id: string;
  readonly username: string;
  readonly stack: string;
  readonly seat_number: number;
  readonly status: PlayerStatus;
  readonly joined_at: string;
};

export type WSPlayerInfo = {
  readonly user_id: string;
  readonly username: string;
  readonly stack: string;
  readonly seat_number: number;
  readonly status: PlayerStatus;
  readonly bet_amount: string;
  readonly is_dealer: boolean;
};

export type WSTableState = {
  readonly table_id: string;
  readonly name: string;
  readonly small_blind: string;
  readonly big_blind: string;
  readonly pot: string;
  readonly community_cards: readonly Card[];
  readonly stage: GameStage;
  readonly dealer_seat: number;
  readonly current_turn: string | null;
  readonly players: readonly WSPlayerInfo[];
};

export type WSCardsDealt = {
  readonly hole_cards: readonly Card[];
  readonly hand_id: string;
};
