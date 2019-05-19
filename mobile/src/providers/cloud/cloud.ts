import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http'; 
import 'rxjs/add/operator/map';

const serverURL = 'http://74.63.216.50:8000';

@Injectable()
export class CloudProvider {
  files:any;

  constructor(private http: HttpClient) {
    this.files = null;
  }

  updateFileList() {
    if(!this.files)
      this.files = [];

    //if (this.files) {
      this.files[0] = serverURL + '/stream.ogg';
      console.log(this.files[0]);
      return this.files;
      //return Promise.resolve(this.files);
    //}
 
    /*return new Promise(resolve => {
      this.http.get(serverURL + '/sound_list')
        .subscribe(data => {
          this.files = data;
          for (let i = 0; i < this.files.length; i++) {
            this.files[i] = serverURL + '/sound?name=' + this.files[i];
          }

          // Shuffle array
          let currentIndex = this.files.length, temporaryValue, randomIndex;
          while (0 !== currentIndex) {
            randomIndex = Math.floor(Math.random() * currentIndex);
            currentIndex -= 1;
            temporaryValue = this.files[currentIndex];
            this.files[currentIndex] = this.files[randomIndex];
            this.files[randomIndex] = temporaryValue;
          }        

          resolve(this.files);
        });
    });*/
  }
}