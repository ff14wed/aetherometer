import { action, computed, observable } from 'mobx';
import { formatDistanceToNow } from 'date-fns';

import Plugin from './plugin';

export default class Stream {
  public id: number;
  public startTime: Date;
  @observable public active: boolean = true;
  @observable public plugins = new Map<string, Plugin>();

  @observable private displayStartTime: string = '';

  constructor(id: number) {
    this.id = id;
    this.startTime = new Date();
    this.updateDisplayStartTime();
  }

  @action public updateDisplayStartTime() {
    this.displayStartTime = formatDistanceToNow(this.startTime, { addSuffix: true });
  }

  @computed public get uniqID() {
    return `${this.id}-${this.startTime.getTime()}`;
  }

  @computed public get displayName() {
    return `Stream ${this.id} - (${this.active ? '' : 'Inactive, '}Started ${this.displayStartTime})`;
  }

  @computed public get shortName() {
    return `Stream ${this.id}`;
  }
}
