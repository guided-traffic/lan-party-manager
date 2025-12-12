import { Injectable, signal, computed } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LatencyService {
  private readonly latencyValues = signal<number[]>([]);
  private readonly maxSamples = 10; // Keep last 10 measurements for averaging

  /** Current average latency in ms */
  readonly latency = computed(() => {
    const values = this.latencyValues();
    if (values.length === 0) return null;
    const sum = values.reduce((a, b) => a + b, 0);
    return Math.round(sum / values.length);
  });

  /** Latency display string */
  readonly latencyDisplay = computed(() => {
    const lat = this.latency();
    if (lat === null) return '...';
    return `${lat}ms`;
  });

  /** Latency status class based on value */
  readonly latencyStatus = computed(() => {
    const lat = this.latency();
    if (lat === null) return 'unknown';
    if (lat < 100) return 'good';
    if (lat < 300) return 'medium';
    return 'slow';
  });

  recordLatency(ms: number): void {
    this.latencyValues.update(values => {
      const newValues = [...values, ms];
      // Keep only the last maxSamples values
      if (newValues.length > this.maxSamples) {
        return newValues.slice(-this.maxSamples);
      }
      return newValues;
    });
  }

  reset(): void {
    this.latencyValues.set([]);
  }
}
