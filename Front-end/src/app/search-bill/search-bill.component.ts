import { Component, OnInit } from '@angular/core';
import { MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { ModalComponent } from '../modal/modal.component';
import { HeaderComponent} from '../header/header.component'
import { APIService } from '../api.service';
import { environment } from 'src/environments/environment';
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

  constructor(  private matdialog: MatDialog) { }

  openAddBillModal(): void{
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;
    dialogConfig.id="modal-component";
    dialogConfig.autoFocus = true;

    dialogConfig.height="350px";
    dialogConfig.width="100%";

    const modalDialog = this.matdialog.open(ModalComponent,dialogConfig)
  }
  ngOnInit(): void {
  }

  SendSearchForm(): void{
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
    let SearchForm = {
      BillId : this.BillId,
      BillStatus : temp,
      Amount : this.Amount,
      ItemId : this.ItemId,
      CustomerId : this.CustomerId,
      CustomerPhone : this.CustomerPhone,
      CustomerName : this.CustomerName,
      ItemName : this.ItemName,
      BillDesc : this.BillDesc
    }
    console.log(SearchForm)
    APIService.POST(`${environment.serverUrl}/v1/Merchant/SearchBill`,SearchForm)
    .then(respone => respone.json())
    .then(respone => respone['SearchBillRespones'])
    .then (data => {
      this.SearchResults = data
      console.log(this.SearchResults)
    })
  }
}
