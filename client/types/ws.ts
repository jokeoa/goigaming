export type WSMessageType =
  | "table_state"
  | "cards_dealt"
  | "player_acted"
  | "community_cards"
  | "hand_result"
  | "player_joined"
  | "player_left"
  | "turn_changed"
  | "new_hand"
  | "error"
  | "pot_updated";

export type WSMessage = {
  readonly type: WSMessageType;
  readonly payload: unknown;
};

export type WSPlayerAction = "fold" | "check" | "call" | "raise" | "all_in";
