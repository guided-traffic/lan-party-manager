import { Injectable, signal } from '@angular/core';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';
import { WebSocketMessage, VotePayload, SettingsPayload } from '../models/websocket.model';
import { Subject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  private connected = signal(false);
  readonly isConnected = this.connected.asReadonly();

  // Subjects for different message types
  readonly voteReceived$ = new Subject<VotePayload>();
  readonly newVote$ = new Subject<VotePayload>();
  readonly settingsUpdate$ = new Subject<SettingsPayload>();

  // General messages observable for timeline component
  private messagesSubject = new Subject<{ type: string; payload: VotePayload }>();
  readonly messages$: Observable<{ type: string; payload: VotePayload }> = this.messagesSubject.asObservable();

  constructor(private authService: AuthService) {}

  connect(): void {
    const token = this.authService.getToken();
    if (!token) {
      console.warn('WebSocket: No token available');
      return;
    }

    // Check if already connected or connecting
    if (this.socket) {
      if (this.socket.readyState === WebSocket.OPEN) {
        console.log('WebSocket: Already connected');
        return;
      }
      if (this.socket.readyState === WebSocket.CONNECTING) {
        console.log('WebSocket: Connection in progress');
        return;
      }
      // Close any existing socket that's closing or closed
      this.socket = null;
    }

    // Clear any pending reconnect timer
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    const wsUrl = `${environment.wsUrl}?token=${token}`;
    console.log('WebSocket: Connecting to', wsUrl);

    try {
      this.socket = new WebSocket(wsUrl);

      this.socket.onopen = () => {
        console.log('WebSocket: Connected successfully');
        this.connected.set(true);
        this.reconnectAttempts = 0;
      };

      this.socket.onmessage = (event) => {
        try {
          const message: WebSocketMessage<VotePayload> = JSON.parse(event.data);
          console.log('WebSocket: Received message', message.type);
          this.handleMessage(message);
        } catch (error) {
          console.error('WebSocket: Failed to parse message', error, event.data);
        }
      };

      this.socket.onclose = (event) => {
        console.log('WebSocket: Disconnected', event.code, event.reason);
        this.connected.set(false);
        this.socket = null;

        // Attempt to reconnect if not a normal closure and user is still authenticated
        if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts && this.authService.isAuthenticated()) {
          this.scheduleReconnect();
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket: Error', error);
      };
    } catch (error) {
      console.error('WebSocket: Failed to create connection', error);
      this.socket = null;
    }
  }

  disconnect(): void {
    // Clear any pending reconnect timer
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }

    if (this.socket) {
      console.log('WebSocket: Disconnecting...');
      this.socket.close(1000, 'User logout');
      this.socket = null;
      this.connected.set(false);
    }
    this.reconnectAttempts = 0;
  }

  private handleMessage(message: WebSocketMessage<VotePayload | SettingsPayload>): void {
    switch (message.type) {
      case 'new_vote':
        console.log('WebSocket: New vote received', message.payload);
        this.newVote$.next(message.payload as VotePayload);
        this.messagesSubject.next({ type: 'new_vote', payload: message.payload as VotePayload });
        break;
      case 'settings_update':
        console.log('WebSocket: Settings update received', message.payload);
        this.settingsUpdate$.next(message.payload as SettingsPayload);
        break;
      default:
        console.log('WebSocket: Unknown message type', message.type);
    }
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * this.reconnectAttempts;
    console.log(`WebSocket: Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      if (this.authService.isAuthenticated()) {
        this.connect();
      }
    }, delay);
  }
}
