import { Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, tap } from 'rxjs';
import { environment } from '../../environments/environment';
import { Settings, UpdateSettingsRequest } from '../models/settings.model';

@Injectable({
  providedIn: 'root'
})
export class SettingsService {
  private settings = signal<Settings | null>(null);
  readonly currentSettings = this.settings.asReadonly();

  constructor(private http: HttpClient) {}

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>(`${environment.apiUrl}/admin/settings`).pipe(
      tap(settings => this.settings.set(settings))
    );
  }

  updateSettings(request: UpdateSettingsRequest): Observable<Settings> {
    return this.http.put<Settings>(`${environment.apiUrl}/admin/settings`, request).pipe(
      tap(settings => this.settings.set(settings))
    );
  }

  // Called by WebSocket service when settings are updated
  applySettingsUpdate(settings: Settings): void {
    this.settings.set(settings);
  }
}
