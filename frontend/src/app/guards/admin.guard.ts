import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

export const adminGuard: CanActivateFn = () => {
  const authService = inject(AuthService);
  const router = inject(Router);

  const user = authService.user();

  if (user?.is_admin) {
    return true;
  }

  // Redirect non-admins to timeline
  router.navigate(['/timeline']);
  return false;
};
