import { Component } from '@angular/core';
import { IonAvatar, IonChip, IonHeader, IonLabel, IonIcon, IonToolbar, IonTitle, IonContent } from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { ellipsisVerticalOutline } from 'ionicons/icons';

@Component({
  selector: 'app-home',
  templateUrl: 'home.page.html',
  styleUrls: ['home.page.scss'],
  standalone: true,
  imports: [IonAvatar, IonChip, IonHeader, IonLabel, IonIcon, IonToolbar, IonTitle, IonContent],
})
export class HomePage {
  constructor() {
    addIcons({ ellipsisVerticalOutline })
  }
}
