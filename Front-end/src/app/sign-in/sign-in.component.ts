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
  
  constructor (private http: Http){
    let username: any;
    let password: any;
  }

  ngOnInit(): void {
  }

  // createAuthorizationHeader(headers: RequestOptionsArgs["headers"]) {
  //   headers.append('Content-Type', 'application/json');
  //   headers.append('Accept', 'application/json');
  //   headers.append('Access-Control-Allow-Origin', 'http://localhost:1234/SignIn');
  //   headers.append('Access-Control-Allow-Credentials', 'true');
  //   headers.append('Authorization', 'Basic ' +
  //     btoa('username:password')); 
  // }

  SendLoginForm(){
    let LoginForm = {
      Username : this.username, 
      Password : this.password,
    }

    console.log(LoginForm)
    // this.http.post(`${environment.serverUrl}`,LoginForm,new RequestOptions({headers: headers})).toPromise()

    APIService.POST(`${environment.serverUrl}/v1/Author`,LoginForm)
  }

}
