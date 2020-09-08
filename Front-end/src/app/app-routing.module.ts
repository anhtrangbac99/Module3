import { APIService } from './api.service';
import { ModalPaymentComponent } from './modal-payment/modal-payment.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { SignInComponent } from './sign-in/sign-in.component';
import { SearchBillComponent } from './search-bill/search-bill.component';
import { ModalComponent } from './modal/modal.component';
const routes: Routes = [
  { path: '', component: SignInComponent },
  { path: 'searchBill', component: SearchBillComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  bootstrap:[
    SignInComponent,
    SearchBillComponent,
    ModalComponent,
    ModalPaymentComponent,
    APIService
  ]
})
export class AppRoutingModule { }
