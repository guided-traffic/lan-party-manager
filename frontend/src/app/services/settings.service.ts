import { Injectable, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';
import { Settings, UpdateSettingsRequest, CreditActionResponse } from '../models/settings.model';

export interface VotingStatusResponse {
  voting_paused: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class SettingsService {
  private settings = signal<Settings | null>(null);
  readonly currentSettings = this.settings.asReadonly();

  // Global voting paused state (can be set via WebSocket before admin loads settings)
  private votingPausedSignal = signal(false);
  readonly votingPaused = this.votingPausedSignal.asReadonly();

  constructor(private http: HttpClient) {
    // Load voting status on service init
    this.loadVotingStatus();
  }

  // Load voting status (accessible to all authenticated users)
  loadVotingStatus(): void {
    this.http.get<VotingStatusResponse>(`${environment.apiUrl}/voting-status`).subscribe({
      next: (response) => {
        this.votingPausedSignal.set(response.voting_paused);
      },
      error: (err) => {
        console.error('Failed to load voting status:', err);
      }
    });
  }

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>(`${environment.apiUrl}/admin/settings`).pipe(
      tap(settings => {
        this.settings.set(settings);
        this.votingPausedSignal.set(settings.voting_paused);
      })
    );
  }

  updateSettings(request: UpdateSettingsRequest): Observable<Settings> {
    return this.http.put<Settings>(`${environment.apiUrl}/admin/settings`, request).pipe(
      tap(settings => {
        this.settings.set(settings);
        this.votingPausedSignal.set(settings.voting_paused);
      })
    );
  }

  resetAllCredits(): Observable<CreditActionResponse> {
    return this.http.post<CreditActionResponse>(`${environment.apiUrl}/admin/credits/reset`, {});
  }

  giveEveryoneCredit(): Observable<CreditActionResponse> {
    return this.http.post<CreditActionResponse>(`${environment.apiUrl}/admin/credits/give`, {});
  }

  // Called by WebSocket service when settings are updated
  applySettingsUpdate(settings: Partial<Settings>): void {
    if (settings.voting_paused !== undefined) {
      this.votingPausedSignal.set(settings.voting_paused);
    }
    const current = this.settings();
    if (current) {
      this.settings.set({ ...current, ...settings });
    }
  }
}
