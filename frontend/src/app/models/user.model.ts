export interface User {
  id: number;
  steam_id: string;
  username: string;
  avatar_url: string;
  avatar_small: string;
  profile_url: string;
}

export interface CurrentUser extends User {
  credits: number;
  seconds_until_credit: number;
  credit_interval_seconds: number;
  credit_max: number;
}
