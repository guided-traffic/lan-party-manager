export interface Game {
  app_id: number;
  name: string;
  header_image_url: string;
  capsule_image_url: string;
  playtime_forever: number;
  categories: string[];
  owner_count: number;
  owners: string[];
  is_pinned: boolean;
}

export interface GamesResponse {
  pinned_games: Game[];
  all_games: Game[];
}
