export type WebSocketMessageType = 'vote_received' | 'new_vote' | 'user_joined' | 'error';

export interface WebSocketMessage<T = unknown> {
  type: WebSocketMessageType;
  payload: T;
}

export interface VotePayload {
  vote_id: number;
  from_user_id: number;
  from_username: string;
  from_avatar: string;
  to_user_id: number;
  to_username: string;
  to_avatar: string;
  achievement_id: string;
  achievement_name: string;
  is_positive: boolean;
  created_at: string;
}
