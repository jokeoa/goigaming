export type WSMessageType =
  | "player_joined"
  | "player_left"
  | "game_started"
  | "card_dealt"
  | "player_action"
  | "community_cards"
  | "pot_updated"
  | "round_ended"
  | "game_ended"
  | "error";

export type WSMessage = {
  readonly type: string;
  readonly game_id: string;
  readonly message: string;
  readonly timestamp: string;
};

export type WSPlayerAction = "fold" | "check" | "call" | "raise" | "all_in";
