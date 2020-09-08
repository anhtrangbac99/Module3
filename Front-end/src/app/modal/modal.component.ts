import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { APIService } from '../api.service';
import { environment } from 'src/environments/environment';
import {Item} from 'src/app/modal/itemclass'

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.css']
})

export class ModalComponent implements OnInit {
  public ListItem : any
  public Message : any
  public Amount : any
  public ItemId : any
  public CustomerId : any
  public CustomerPhone : any
  public BillDesc : any
  public CustomerName : any
  public ItemName : any
  public IsCustomerExist = false

  public ListItemCreate: Array<Item> = []
  constructor(public dialogRef: MatDialogRef<ModalComponent>) { }

  ngOnInit(): void {
    APIService.CheckToken('/')
    APIService.GET(`${environment.serverUrl}/v1/Merchant/ListItem`)
    .then (respone => respone.json())
    .then (respone => respone['item'])
    .then(
      data => {
        this.ListItem = data
        for (var value of this.ListItem){
          let item = new Item()
          item.ItemId = value['ItemId']
          item.Amount = 0
          this.ListItemCreate.push(item)
        }
        console.log(this.ListItem)        
        console.log(this.ListItemCreate)

      }
    )
  }

  ChangeAmount(index,event) : void {
    this.ListItemCreate[index].Amount = event.target.value
    console.log(this.ListItemCreate[index])
  }
  addBill():void {
    APIService.CheckToken('/')

    var addBillForm = {
      Item : this.ListItemCreate,
      CustomerId: this.CustomerId,
      BillDesc : this.BillDesc
    }
    console.log(addBillForm)
    if (this.IsCustomerExist){
      APIService.POST(`${environment.serverUrl}/v1/Merchant/`+localStorage.getItem('UserToken')+`/CreateBill`,addBillForm)
      .then(respone => respone.json())
      .then(data => console.log(data))
      this.dialogRef.close()
      window.location.href = '/searchBill'
    }
    
    this.Message = "Customer does not exist"

    /*************** End Here ***************/
    // window.location.href = '/searchBill'
  }

  getCustomer():void{
    APIService.CheckToken('/')

    var Json = {
      CustomerPhone: this.CustomerPhone
    }
    APIService.GET(`${environment.serverUrl}/v1/Merchant/Customer/`+ this.CustomerPhone )
    .then (respone => respone.json())
    .then(
      data => {
        if (data['CustomerName']!=null){
          this.IsCustomerExist = true
        }
        else {
          this.IsCustomerExist = false
        }
        console.log(this.IsCustomerExist)
        this.CustomerName = data['CustomerName']
        this.CustomerId = data['CustomerId']
      }
    )
  }

  // setItemName(ItemName : string): void {
  //   console.log(ItemName)
  //   this.CustomerName = ItemName
  // }
  closeModal(): void{
    APIService.CheckToken('/')
    this.dialogRef.close()
  }
}
