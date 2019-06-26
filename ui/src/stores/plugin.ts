import { action, computed, observable } from 'mobx';

interface PluginParams {
  apiURL: string;
  apiToken: string;
  streamID: number;
}

export default class Plugin {
  @observable public name: string;
  @observable public url: string;
  @observable public params?: PluginParams;
  @observable public icon?: string;

  private streamUniqID: string;

  constructor(
    name: string,
    url: string,
    streamUniqID?: string,
    params?: PluginParams,
  ) {
    this.name = name;
    this.url = url;
    this.params = params;
    this.streamUniqID = streamUniqID || '';
  }

  @action public setIcon(icon?: string) {
    this.icon = icon;
  }

  public withParams(streamUniqID: string, params: PluginParams) {
    return new Plugin(this.name, this.url, streamUniqID, params);
  }

  @computed public get id() {
    return `plugin-${this.streamUniqID}-${this.name}`;
  }
}
