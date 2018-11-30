import { Injectable } from '@angular/core';
import { of } from 'rxjs';

@Injectable()
export class CloudProvider {
  files:any = [
    { url: 'http://51.15.88.55:3000/sound', 
      name: 'Sound by Taj'
    }
  ];
  getFiles() {
   return of(this.files);
  }
}