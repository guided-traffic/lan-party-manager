import { Injectable, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';
import { CurrentUser } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly TOKEN_KEY = 'lan_party_token';

  private currentUser = signal<CurrentUser | null>(null);
  private loading = signal(false);

  readonly user = this.currentUser.asReadonly();
  readonly isAuthenticated = computed(() => !!this.currentUser());
  readonly isLoading = this.loading.asReadonly();
  readonly credits = computed(() => this.currentUser()?.credits ?? 0);

  constructor(
    private http: HttpClient,
    private router: Router
  ) {
    // Check for existing token on startup
    if (this.getToken()) {
      this.loadCurrentUser();
    }
  }

  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  removeToken(): void {
    localStorage.removeItem(this.TOKEN_KEY);
  }

  login(): void {
    // Redirect to backend Steam auth endpoint
    window.location.href = `${environment.apiUrl}/auth/steam`;
  }

  handleCallback(token: string): void {
    this.setToken(token);
    this.loadCurrentUser();
  }

  logout(): void {
    this.http.post(`${environment.apiUrl}/auth/logout`, {}).subscribe({
      complete: () => {
        this.removeToken();
        this.currentUser.set(null);
        this.router.navigate(['/login']);
      },
      error: () => {
        // Even if the API call fails, clear local state
        this.removeToken();
        this.currentUser.set(null);
        this.router.navigate(['/login']);
      }
    });
  }

  loadCurrentUser(): void {
    this.loading.set(true);
    this.http.get<{ user: CurrentUser }>(`${environment.apiUrl}/auth/me`).subscribe({
      next: (response) => {
        this.currentUser.set(response.user);
        this.loading.set(false);
      },
      error: (error) => {
        console.error('Failed to load user:', error);
        this.removeToken();
        this.currentUser.set(null);
        this.loading.set(false);
        this.router.navigate(['/login']);
      }
    });
  }

  checkAuth(): void {
    if (this.getToken()) {
      this.loadCurrentUser();
    }
  }

  updateCredits(credits: number): void {
    const user = this.currentUser();
    if (user) {
      this.currentUser.set({ ...user, credits });
    }
  }

  refreshUser(): void {
    this.loadCurrentUser();
  }
}
