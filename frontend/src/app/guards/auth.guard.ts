import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  // getToken() automatically returns null for expired tokens
  const token = authService.getToken();
  const hasValidToken = !!token;
  console.log('[AuthGuard] Checking access - Has valid token:', hasValidToken);

  if (hasValidToken) {
    return true;
  }

  console.log('[AuthGuard] No valid token, redirecting to login');
  router.navigate(['/login']);
  return false;
};
