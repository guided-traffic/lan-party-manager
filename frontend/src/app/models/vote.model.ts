import { Achievement } from './achievement.model';
import { User } from './user.model';

export interface Vote {
  id: number;
  from_user: User;
  to_user: User;
  achievement_id: string;
  achievement: Achievement;
  created_at: string;
}

export interface CreateVoteRequest {
  to_user_id: number;
  achievement_id: string;
}

export interface VoteResponse {
  vote: Vote;
  credits: number;
}

export interface LeaderboardEntry {
  user: User;
  vote_count: number;
  rank: number;
}

export interface AchievementLeaderboard {
  achievement: Achievement;
  leaders: LeaderboardEntry[];
}
