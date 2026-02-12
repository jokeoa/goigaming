export type Suit = "hearts" | "diamonds" | "clubs" | "spades";

export type Rank =
  | "2"
  | "3"
  | "4"
  | "5"
  | "6"
  | "7"
  | "8"
  | "9"
  | "10"
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
  | "showdown";

export type PlayerInfo = {
  readonly id: string;
  readonly username: string;
  readonly seatIndex: number;
  readonly chips: number;
  readonly currentBet: number;
  readonly isFolded: boolean;
  readonly isDealer: boolean;
  readonly isActive: boolean;
  readonly cards?: readonly Card[];
};

export type PokerTable = {
  readonly id: string;
  readonly name: string;
  readonly maxPlayers: number;
  readonly currentPlayers: number;
  readonly smallBlind: number;
  readonly bigBlind: number;
  readonly minBuyIn: number;
  readonly maxBuyIn: number;
  readonly stage: GameStage;
};

export type TableState = {
  readonly table: PokerTable;
  readonly players: readonly PlayerInfo[];
  readonly communityCards: readonly Card[];
  readonly pot: number;
  readonly currentTurn: string | null;
  readonly stage: GameStage;
};
