export interface Settings {
  credit_interval_minutes: number;
  credit_max: number;
  voting_paused: boolean;
}

export interface UpdateSettingsRequest {
  credit_interval_minutes?: number;
  credit_max?: number;
  voting_paused?: boolean;
}

export interface CreditActionResponse {
  message: string;
  users_affected: number;
}
