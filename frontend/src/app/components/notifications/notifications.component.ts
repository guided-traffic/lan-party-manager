import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { NotificationService, Notification } from '../../services/notification.service';

@Component({
  selector: 'app-notifications',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="notifications-container">
      @for (notification of notifications.all(); track notification.id) {
        <div
          class="notification"
          [class]="'notification-' + notification.type"
          [class.positive]="notification.isPositive"
          [class.negative]="notification.isPositive === false"
        >
          @if (notification.avatar) {
            <img [src]="notification.avatar" class="notification-avatar" alt="" />
          }
          <div class="notification-content">
            <div class="notification-title">{{ notification.title }}</div>
            <div class="notification-message">{{ notification.message }}</div>
          </div>
          <button class="notification-close" (click)="dismiss(notification.id)">×</button>
        </div>
      }
    </div>
  `,
  styles: [`
    @use 'variables' as *;

    .notifications-container {
      position: fixed;
      top: 80px;
      right: 24px;
      z-index: 1000;
      display: flex;
      flex-direction: column;
      gap: 12px;
      max-width: 380px;
    }

    .notification {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      padding: 16px;
      background: $bg-card;
      border: 1px solid $border-color;
      border-radius: $radius-lg;
      box-shadow: $shadow-lg;
      animation: slideIn 0.3s ease;

      &-vote {
        &.positive {
          border-color: rgba($accent-positive, 0.4);
          background: linear-gradient(135deg, rgba($accent-positive, 0.1), transparent);
        }

        &.negative {
          border-color: rgba($accent-negative, 0.4);
          background: linear-gradient(135deg, rgba($accent-negative, 0.1), transparent);
        }
      }

      &-success {
        border-color: rgba($accent-success, 0.4);
      }

      &-error {
        border-color: rgba($accent-error, 0.4);
      }

      &-info {
        border-color: rgba($accent-primary, 0.4);
      }
    }

    .notification-avatar {
      width: 40px;
      height: 40px;
      border-radius: 50%;
      flex-shrink: 0;
    }

    .notification-content {
      flex: 1;
    }

    .notification-title {
      font-weight: 600;
      font-size: 14px;
      margin-bottom: 4px;
    }

    .notification-message {
      font-size: 13px;
      color: $text-secondary;
      line-height: 1.4;
    }

    .notification-close {
      background: none;
      border: none;
      color: $text-muted;
      font-size: 20px;
      cursor: pointer;
      padding: 0;
      line-height: 1;
      transition: color $transition-fast;

      &:hover {
        color: $text-primary;
      }
    }

    @keyframes slideIn {
      from {
        opacity: 0;
        transform: translateX(100px);
      }
      to {
        opacity: 1;
        transform: translateX(0);
      }
    }
  `]
})
export class NotificationsComponent {
  notifications = inject(NotificationService);

  dismiss(id: number): void {
    this.notifications.dismiss(id);
  }
}
