export interface Settings {
  credit_interval_minutes: number;
  credit_max: number;
}

export interface UpdateSettingsRequest {
  credit_interval_minutes?: number;
  credit_max?: number;
}
