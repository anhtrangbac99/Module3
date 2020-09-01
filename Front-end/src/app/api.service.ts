import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';

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
