import { Component, inject } from '@angular/core';
import { VersionService } from '../../services/version.service';
import { WebSocketService } from '../../services/websocket.service';
import { LatencyService } from '../../services/latency.service';

@Component({
  selector: 'app-footer',
  standalone: true,
  templateUrl: './footer.html',
  styleUrl: './footer.scss'
})
export class FooterComponent {
  currentYear = new Date().getFullYear();
  versionService = inject(VersionService);
  ws = inject(WebSocketService);
  latencyService = inject(LatencyService);
}
