import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { delay } from 'rxjs';
import { Observable } from 'rxjs';

import { ApiResult, EntryResult } from './interfaces';

const apiUrl = '';

@Injectable({
  providedIn: 'root'
})
export class EntryService {
  private http = inject(HttpClient);

  constructor() { }

  getTopEntries(page = 1): Observable<ApiResult> {
    return this.http
      .get<ApiResult>(apiUrl)
      .pipe(
        delay(2000) // simulating a slow network here to test
      )
  }

  getEntryDetails(id: string): Observable<EntryResult> {
    return this.http.get<EntryResult>(`${apiUrl}/${id}`)
  }
}
