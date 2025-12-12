import { HttpInterceptorFn, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, tap, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { ConnectionStatusService } from '../services/connection-status.service';
import { LatencyService } from '../services/latency.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);
  const connectionStatus = inject(ConnectionStatusService);
  const latencyService = inject(LatencyService);

  const token = authService.getToken();
  const startTime = performance.now();

  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }

  return next(req).pipe(
    tap(event => {
      if (event instanceof HttpResponse) {
        const latency = Math.round(performance.now() - startTime);
        latencyService.recordLatency(latency);
      }
    }),
    catchError((error: HttpErrorResponse) => {
      console.log('[AuthInterceptor] Error:', error.status, req.url);

      if (error.status === 401) {
        console.log('[AuthInterceptor] 401 - Removing token and redirecting to login');
        authService.removeToken();
        router.navigate(['/login']);
      } else if (error.status === 0 || error.status >= 500) {
        // Network error or server error - backend is unavailable
        console.log('[AuthInterceptor] Backend unavailable:', error.status);
        connectionStatus.setDisconnected();
      }

      return throwError(() => error);
    })
  );
};
