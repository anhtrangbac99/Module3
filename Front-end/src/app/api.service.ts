import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class APIService {
  constructor() { }

  static POST(url,data) {
    return fetch(url, {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/json'
      }),
      body : data
    }).then(response => response.json())
    .then (data => console.log(data['IsExisted']))
  }
}
