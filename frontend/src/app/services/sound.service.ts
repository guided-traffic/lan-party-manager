import { Injectable } from '@angular/core';

export type SoundType = 'new-credit' | 'good-review' | 'bad-review' | 'review-given' | 'new-king';

@Injectable({
  providedIn: 'root'
})
export class SoundService {
  private audioCache = new Map<SoundType, HTMLAudioElement>();
  private soundEnabled = true;

  constructor() {
    // Preload sounds
    this.preloadSound('new-credit');
    this.preloadSound('good-review');
    this.preloadSound('bad-review');
    this.preloadSound('review-given');
    this.preloadSound('new-king');
  }

  private preloadSound(type: SoundType): void {
    const audio = new Audio(`/sounds/${type}.mp3`);
    audio.preload = 'auto';
    this.audioCache.set(type, audio);
  }

  play(type: SoundType): void {
    if (!this.soundEnabled) {
      return;
    }

    // Clone the audio element to allow overlapping sounds
    const cachedAudio = this.audioCache.get(type);
    if (cachedAudio) {
      const audio = cachedAudio.cloneNode() as HTMLAudioElement;
      audio.volume = 0.5;
      audio.play().catch(err => {
        // Browser may block autoplay, that's okay
        console.warn('Sound playback blocked:', err);
      });
    }
  }

  playNewCredit(): void {
    this.play('new-credit');
  }

  playGoodReview(): void {
    this.play('good-review');
  }

  playBadReview(): void {
    this.play('bad-review');
  }

  playReviewGiven(): void {
    this.play('review-given');
  }

  playNewKing(): void {
    this.play('new-king');
  }

  setSoundEnabled(enabled: boolean): void {
    this.soundEnabled = enabled;
  }

  isSoundEnabled(): boolean {
    return this.soundEnabled;
  }
}
