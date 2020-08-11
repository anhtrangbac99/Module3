import { Component, OnInit } from '@angular/core';
import { MatDialog, MatDialogConfig, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.css']
})
export class ModalComponent implements OnInit {

  constructor(public dialogRef: MatDialogRef<ModalComponent>) { }

  ngOnInit(): void {
  }

  addBill():void {
    /*************** Add Bill Start Here ***************/



    /*************** End Here ***************/

    this.closeModal();
  }

  closeModal(): void{
    this.dialogRef.close();
  }
}
