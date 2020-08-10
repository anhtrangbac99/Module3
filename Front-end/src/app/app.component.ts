import { Component } from '@angular/core';
import {Router} from '@angular/router';


@Component({
  selector: 'app-root',
  template: '<router-outlet></router-outlet>',
  styles: []
})
export class AppComponent {
  title = 'Angular Front-end';

  constructor(private router:Router){}

  GoToPage(PageName:string):void{
    this.router.navigate(['${PageName}']);
  }
}


