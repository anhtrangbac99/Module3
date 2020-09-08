import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import 'whatwg-fetch'

@Injectable({
  providedIn: 'root'
})
export class APIService {
  static http: HttpClient;
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
        'Content-Type': 'application/json',
      }),
    })
  }

  // static GET(url){
  //   return this.http.get(url)
  // }

  static CheckToken(url:string) : void{
    let UserToken = localStorage.getItem('UserToken')
    if (UserToken){
      APIService.GET(`${environment.serverUrl}/v1/Merchant/UserToken/` + UserToken)
      .then (respone => respone.json())
      .then (data => {
        if (data['IsExisted'] == -1) {
          window.location.href= url
        }
      })
    }
    else {
      window.location.href= url
    }
  }

}
