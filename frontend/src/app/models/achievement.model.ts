export interface Achievement {
  id: string;
  name: string;
  description: string;
  image_url: string;
  is_positive: boolean;
}

export interface AchievementsResponse {
  achievements: Achievement[];
  positive: Achievement[];
  negative: Achievement[];
}
