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
      body : JSON.stringify(data)
    })
  }

  static GET(url){
    return fetch(url, {
      method: 'GET',
      headers: new Headers({
        'Content-Type': 'application/json'
      }),
      //body : JSON.stringify(data)
    })
  }
}
