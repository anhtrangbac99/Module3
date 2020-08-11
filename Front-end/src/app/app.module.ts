import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';

import { SignInComponent } from './sign-in/sign-in.component';
import { SearchBillComponent } from './search-bill/search-bill.component';
import { HeaderComponent } from './header/header.component';
import { ModalComponent } from './modal/modal.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

import {MatButtonModule} from '@angular/material/button';
import {MatDialogModule} from '@angular/material/dialog' ;

@NgModule({
  declarations: [
    AppComponent,
    SignInComponent,
    SearchBillComponent,
    HeaderComponent,
    ModalComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatButtonModule,
    MatDialogModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
