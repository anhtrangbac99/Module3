import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';
import { APIService } from '../api.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-modal',
  templateUrl: './modal-payment.component.html',
  styleUrls: ['./modal-payment.component.css']
})
export class ModalPaymentComponent implements OnInit {
  public BillDetailForm : any
  constructor(public dialogRef: MatDialogRef<ModalPaymentComponent>) { }

  ngOnInit(): void {
    let BillId = localStorage.getItem("BillId")

    APIService.GET(`${environment.serverUrl}/v1/Merchant/BillId/` + BillId)
    .then (respone => respone.json())
    .then (data =>{
      this.BillDetailForm = data
    })
  }

  closeModal(): void{
    localStorage.removeItem('BillId')
    this.dialogRef.close()
  }
}
