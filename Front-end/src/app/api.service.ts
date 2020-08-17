import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class APIService {
  constructor() { }

  static POST(url,data) {
    let headers = new Headers()
    headers.append('Content-Type', 'application/json');
    headers.append('Accept', 'application/json');

    headers.append('Access-Control-Allow-Origin', url);
    headers.append('Access-Control-Allow-Credentials', 'true');

    return fetch(url, {
      mode: 'no-cors',
      credentials: 'include',
      method: 'POST',
      headers: headers,
      body : JSON.stringify(data)
    }).then(function(response) {
      console.log(response.status)
      if (!response.ok) {
        throw new Error('Bad status code from server.');
      }
      return response.json();
    })
    .then(data => console.log(data));

  }
}
