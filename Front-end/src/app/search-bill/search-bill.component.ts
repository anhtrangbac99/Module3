import { Component, OnInit } from '@angular/core';
import { MatDialogConfig, MatDialog } from '@angular/material/dialog';
import { ModalComponent } from '../modal/modal.component';
import { HeaderComponent} from '../header/header.component'
@Component({
  selector: 'app-search-bill',
  templateUrl: './search-bill.component.html',
  styleUrls: ['./search-bill.component.css']
})
export class SearchBillComponent implements OnInit {

  constructor(public matdialog:MatDialog) { }

  openAddBillModal(): void{
    const dialogConfig = new MatDialogConfig();
    dialogConfig.disableClose = true;
    dialogConfig.id="modal-component";
    dialogConfig.height="350px";
    dialogConfig.width="600px";

    const modalDialog = this.matdialog.open(ModalComponent,dialogConfig)
  }
  ngOnInit(): void {
  }

}
