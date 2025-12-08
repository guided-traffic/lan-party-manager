import { Injectable, signal } from '@angular/core';

export interface Notification {
  id: number;
  type: 'success' | 'error' | 'info' | 'vote';
  title: string;
  message: string;
  avatar?: string;
  isPositive?: boolean;
  duration?: number;
}

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private notifications = signal<Notification[]>([]);
  private nextId = 1;

  readonly all = this.notifications.asReadonly();

  show(notification: Omit<Notification, 'id'>): void {
    const id = this.nextId++;
    const duration = notification.duration ?? 5000;

    this.notifications.update(list => [...list, { ...notification, id }]);

    if (duration > 0) {
      setTimeout(() => this.dismiss(id), duration);
    }
  }

  dismiss(id: number): void {
    this.notifications.update(list => list.filter(n => n.id !== id));
  }

  success(title: string, message: string): void {
    this.show({ type: 'success', title, message });
  }

  error(title: string, message: string): void {
    this.show({ type: 'error', title, message, duration: 8000 });
  }

  info(title: string, message: string): void {
    this.show({ type: 'info', title, message });
  }

  voteReceived(fromUsername: string, achievementName: string, avatar: string, isPositive: boolean): void {
    this.show({
      type: 'vote',
      title: isPositive ? '🎉 Achievement erhalten!' : '😅 Achievement erhalten',
      message: `${fromUsername} hat dir "${achievementName}" gegeben!`,
      avatar,
      isPositive,
      duration: 8000
    });
  }
}
