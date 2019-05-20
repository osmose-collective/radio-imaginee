import {Component, ViewChild} from '@angular/core';
import {trigger, state, style, animate, transition } from '@angular/animations';
import {NavController, ModalController, NavParams, Navbar, Content, LoadingController, MenuController} from 'ionic-angular';
import {Store} from '@ngrx/store';
import { InAppBrowser } from '@ionic-native/in-app-browser'

import {CANPLAY, PLAYING, LOADSTART, RESET} from '../../providers/store/store';
import {AudioProvider} from '../../providers/audio/audio';

const radioURL = 'http://stream.osmose.world/radio-imaginee-192.mp3';

@Component({
  selector: 'page-home',
  templateUrl: 'home.html',
  animations: [
    trigger('showHide', [
      state(
        'active',
        style({
          opacity: 1
        })
      ),
      state(
        'inactive',
        style({
          opacity: 0
        })
      ),
      transition('inactive => active', animate('250ms ease-in')),
      transition('active => inactive', animate('250ms ease-out'))
    ])
  ]
})

export class HomePage {
  state: any = {};
  onSeekState: boolean;
  toggleMenu: boolean;
  @ViewChild(Navbar) navBar: Navbar;
  @ViewChild(Content) content: Content;

  constructor(
    public navCtrl: NavController,
    public modalCtrl: ModalController,
    public menuCtrl: MenuController,
    public navParams: NavParams,
    public audioProvider: AudioProvider,
    public loadingCtrl: LoadingController,
    private iab: InAppBrowser,
    private store: Store<any>
  ) {
    this.getDocuments();
  }

  getDocuments() {
    let loader = this.presentLoading();
    loader.dismiss();
    this.initStream(radioURL);
    this.toggleMenu = false;
  }

  presentLoading() {
    let loading = this.loadingCtrl.create({
      content: 'Loading Content. Please Wait...'
    });
    loading.present();
    return loading;
  }

  ionViewWillLoad() {
    this.store.select('appState').subscribe((value: any) => {
      this.state = value.media;
    });
  }

  resetState() {
    this.audioProvider.stop();
    this.store.dispatch({ type: RESET });
  }

  initStream(url) {
    this.resetState();
    this.audioProvider.initStream(url).subscribe(event => {

      switch (event.type) {
        case 'canplay':
          return this.store.dispatch({ type: CANPLAY, payload: { value: true } });

        case 'playing':
          return this.store.dispatch({ type: PLAYING, payload: { value: true } });

        case 'pause':
          return this.store.dispatch({ type: PLAYING, payload: { value: false } });

        case 'loadstart':
          return this.store.dispatch({ type: LOADSTART, payload: { value: true } });
      }
    });
  }

  pause() {
    this.audioProvider.pause();
  }

  play() {
    this.audioProvider.play();
  }

  stop() {
    this.audioProvider.stop();
  }

  onSeekStart() {
    this.onSeekState = this.state.playing;
    if (this.onSeekState) {
      this.pause();
    }
  }

  onSeekEnd(event) {
    if (this.onSeekState) {
      this.audioProvider.seekTo(event.value);
      this.play();
    } else {
      this.audioProvider.seekTo(event.value);
    }
  }

  openLaSuiteDuMondeLink() {
    this.iab.create('https://www.lasuitedumonde.com/', '_system');
  }

  openLaSuiteDuMondeInstagramLink()
  {
    this.iab.create('https://www.instagram.com/lasuitedumonde/', '_system');
  }

  openLaSuiteDuMondeFacebookLink()
  {
    this.iab.create('https://www.facebook.com/CommuneDuBandiat/', '_system');
  }

  openLaSuiteDuMondeTelegramLink(){
    this.iab.create('https://t.me/joinchat/FQtyO1TnGJvJ2qqrpYHhnQ', '_system');
  }

  openOSMOSELink(){
    this.iab.create('https://www.osmosecollective.com', '_system');
  }

  reset() {
    this.resetState();
  }
}