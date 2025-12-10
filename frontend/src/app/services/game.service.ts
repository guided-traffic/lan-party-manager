import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { GamesResponse } from '../models/game.model';

@Injectable({
  providedIn: 'root'
})
export class GameService {
  constructor(private http: HttpClient) {}

  getMultiplayerGames(): Observable<GamesResponse> {
    return this.http.get<GamesResponse>(`${environment.apiUrl}/games`);
  }

  refreshGames(): Observable<GamesResponse> {
    return this.http.post<GamesResponse>(`${environment.apiUrl}/games/refresh`, {});
  }
}
