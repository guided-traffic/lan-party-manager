import { Injectable, signal } from '@angular/core';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';
import { WebSocketMessage, VotePayload } from '../models/websocket.model';
import { Subject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private socket: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;

  private connected = signal(false);
  readonly isConnected = this.connected.asReadonly();

  // Subjects for different message types
  readonly voteReceived$ = new Subject<VotePayload>();
  readonly newVote$ = new Subject<VotePayload>();

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

    if (this.socket?.readyState === WebSocket.OPEN) {
      console.log('WebSocket: Already connected');
      return;
    }

    const wsUrl = `${environment.wsUrl}?token=${token}`;
    console.log('WebSocket: Connecting...');

    this.socket = new WebSocket(wsUrl);

    this.socket.onopen = () => {
      console.log('WebSocket: Connected');
      this.connected.set(true);
      this.reconnectAttempts = 0;
    };

    this.socket.onmessage = (event) => {
      try {
        const message: WebSocketMessage<VotePayload> = JSON.parse(event.data);
        this.handleMessage(message);
      } catch (error) {
        console.error('WebSocket: Failed to parse message', error);
      }
    };

    this.socket.onclose = (event) => {
      console.log('WebSocket: Disconnected', event.code, event.reason);
      this.connected.set(false);
      this.socket = null;

      // Attempt to reconnect if not a normal closure
      if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
        this.scheduleReconnect();
      }
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket: Error', error);
    };
  }

  disconnect(): void {
    if (this.socket) {
      this.socket.close(1000, 'User logout');
      this.socket = null;
      this.connected.set(false);
    }
  }

  private handleMessage(message: WebSocketMessage<VotePayload>): void {
    switch (message.type) {
      case 'vote_received':
        console.log('WebSocket: Vote received notification', message.payload);
        this.voteReceived$.next(message.payload);
        this.messagesSubject.next({ type: 'vote_received', payload: message.payload });
        break;
      case 'new_vote':
        console.log('WebSocket: New vote in timeline', message.payload);
        this.newVote$.next(message.payload);
        this.messagesSubject.next({ type: 'new_vote', payload: message.payload });
        break;
      default:
        console.log('WebSocket: Unknown message type', message.type);
    }
  }

  private scheduleReconnect(): void {
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * this.reconnectAttempts;
    console.log(`WebSocket: Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);

    setTimeout(() => {
      if (this.authService.isAuthenticated()) {
        this.connect();
      }
    }, delay);
  }
}
