import { Component, OnInit, Type } from '@angular/core';
import { MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { ModalComponent } from '../modal/modal.component';
import { HeaderComponent} from '../header/header.component'
import { APIService } from '../api.service';
import { environment } from 'src/environments/environment';
import { flatMap } from 'rxjs/operators';
import { ModalPaymentComponent } from '../modal-payment/modal-payment.component';
@Component({
  selector: 'app-search-bill',
  templateUrl: './search-bill.component.html',
  styleUrls: ['./search-bill.component.css']
})
export class SearchBillComponent implements OnInit {
  public BillId : any  
  public BillStatus : any
  public Amount : any
  public ItemId : any
  public CustomerId : any
  public CustomerPhone : any
  public BillDesc : any
  public CustomerName : any
  public ItemName : any

  public SearchResults : any
  public Authorized : any
  constructor(  private matdialog: MatDialog) { }

  openAddBillModal(): void{
    APIService.CheckToken('/')

    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;
    dialogConfig.id="modal-component";
    dialogConfig.autoFocus = true;

    dialogConfig.height="350px";
    dialogConfig.width="100%";

    const modalDialog = this.matdialog.open(ModalComponent,dialogConfig)
  }

  openPayModal(billId:any): void{
    APIService.CheckToken('/')

    localStorage.setItem('BillId',billId)
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;
    dialogConfig.id="modal-component";
    dialogConfig.autoFocus = true;

    dialogConfig.height="350px";
    dialogConfig.width="50%";

    const modalDialog = this.matdialog.open(ModalPaymentComponent,dialogConfig)
  }

  ngOnInit(): void {
    let UserToken = localStorage.getItem('UserToken')
    if (UserToken){
      APIService.GET(`${environment.serverUrl}/v1/Merchant/UserToken/` + UserToken)
      .then (respone => respone.json())
      .then (data => {
        if (data['IsExisted'] == -1) {
          window.location.href = '/'
          return
        }
        this.Authorized = data['Authorized']
      })

    } 
    else {
      window.location.href = '/'
    }
  }

  SendSearch(): void {
    let UserToken = localStorage.getItem('UserToken')
    var str : string

    APIService.CheckToken('/')
    let temp : any
    if (this.BillStatus == 'Mới tạo'){
      temp = 1
    }
    if (this.BillStatus == 'Đã thanh toán'){
      temp = 2
    }
    if (this.BillStatus == 'Đã hủy'){
      temp = 3
    }
    if (!this.BillId){
      this.BillId = 0
    }
    if (!temp){
      temp = 0
    }
    if (!this.Amount){
      this.Amount = 0
    }
    if (!this.ItemId){
      this.ItemId = 0
    }
    if (!this.CustomerId){
      this.CustomerId = 0
    }
    if (!this.CustomerPhone){
      this.CustomerPhone = " "
    }
    if (!this.CustomerName){
      this.CustomerName = " "
    }
    if (!this.ItemName){
      this.ItemName = " "
    }
    if (!this.BillDesc){
      this.BillDesc = " "
    }

    APIService.GET(`${environment.serverUrl}/v1/Merchant/Search/` + UserToken + `/BillId/` + this.BillId + `/BillStatus/`+temp+`/Amount/`+this.Amount+`/ItemId/`+this.ItemId+`/CustomerId/`+ this.CustomerId+`/CustomerPhone/`+this.CustomerPhone+`/CustomerName/`+this.CustomerName+`/ItemName/`+this.ItemName + `/BillDesc/`+this.BillDesc)
    .then(respone => respone.json())
    .then(respone => respone['SearchRespones'])
    .then (data => {
      this.SearchResults = data
      console.log(this.SearchResults)
    })
  }
    
}
