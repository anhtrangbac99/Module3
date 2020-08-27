import { APIService } from './../api.service';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit, Injectable } from '@angular/core';
import {Http, RequestOptions, RequestMethod, RequestOptionsArgs} from '@angular/http';
import { environment } from 'src/environments/environment';
import axios from 'axios';
import { throwError } from 'rxjs/internal/observable/throwError';
import { map, catchError } from 'rxjs/operators';
import {Headers} from '@angular/http'
@Component({
  selector: 'Sign-in',
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})

export class SignInComponent implements OnInit {
  public username :any
  public password : any
  public message = " "
  
  constructor (private http: Http){
    let username: any;
    let password: any;
  }

  ngOnInit(): void {
  }


  SendLoginForm(){
    let LoginForm = {
      Username : this.username, 
      Password : this.password,
    }

    console.log(LoginForm)
    // this.http.post(`${environment.serverUrl}`,LoginForm,new RequestOptions({headers: headers})).toPromise()

    APIService.POST(`${environment.serverUrl}/v1/Merchant/Author`,LoginForm)
    .then(response => response.json())
    .then(data => {
      if (data['IsExisted'] == 1){
        if (data['User_Id'] != -1){
          if (data['Authorized']==1){
            window.location.href = '/searchBill'
            return
          }
        }
      }
      this.message = "Username or password is incorrect"
    })
  }

}
