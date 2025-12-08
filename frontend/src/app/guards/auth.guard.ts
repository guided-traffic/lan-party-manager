import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const authGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const hasToken = !!authService.getToken();
  console.log('[AuthGuard] Checking access - Has token:', hasToken);

  if (hasToken) {
    return true;
  }

  console.log('[AuthGuard] No token, redirecting to login');
  router.navigate(['/login']);
  return false;
};
