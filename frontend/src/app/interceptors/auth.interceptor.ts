import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const token = authService.getToken();

  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }

  return next(req).pipe(
    catchError((error: HttpErrorResponse) => {
      console.log('[AuthInterceptor] Error:', error.status, req.url);
      if (error.status === 401) {
        console.log('[AuthInterceptor] 401 - Removing token and redirecting to login');
        authService.removeToken();
        router.navigate(['/login']);
      }
      return throwError(() => error);
    })
  );
};
