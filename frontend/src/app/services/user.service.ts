import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../../environments/environment';
import { User } from '../models/user.model';

@Injectable({
  providedIn: 'root'
})
export class UserService {
  constructor(private http: HttpClient) {}

  getAll(): Observable<User[]> {
    return this.http.get<{ users: User[] }>(`${environment.apiUrl}/users`)
      .pipe(map(response => response.users));
  }

  getOthers(): Observable<User[]> {
    return this.http.get<{ users: User[] }>(`${environment.apiUrl}/users/others`)
      .pipe(map(response => response.users));
  }

  getById(id: number): Observable<User> {
    return this.http.get<{ user: User }>(`${environment.apiUrl}/users/${id}`)
      .pipe(map(response => response.user));
  }
}
