import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../environments/environment';
import { Achievement, AchievementsResponse } from '../models/achievement.model';

@Injectable({
  providedIn: 'root'
})
export class AchievementService {
  constructor(private http: HttpClient) {}

  getAll(): Observable<AchievementsResponse> {
    return this.http.get<AchievementsResponse>(`${environment.apiUrl}/achievements`);
  }

  getById(id: string): Observable<Achievement> {
    return this.http.get<{ achievement: Achievement }>(`${environment.apiUrl}/achievements/${id}`)
      .pipe(map(response => response.achievement));
  }
}
