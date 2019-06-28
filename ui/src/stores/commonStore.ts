import { action, computed, observable, autorun, IReactionDisposer } from 'mobx';
import GQLClient from '../api/gqlClient';
import { ipcRenderer, IpcMessageEvent } from 'electron';

import Plugin from './plugin';
import Stream from './stream';

import Store from 'electron-store';

import { JSONSchema } from 'json-schema-typed';
import { ZenObservable } from 'zen-observable-ts';

interface RendererPayload {
  apiPort: number;
  adminOTP?: string;
  adminToken?: string;
}
interface TreeViewSpec {
  id: string;
  name: string;
  children?: TreeViewSpec[];
  [k: string]: any;
}

const schema: { [key: string]: JSONSchema } = {
  defaultPlugins: {
    type: 'array',
    items: {
      required: ['name', 'url'],
      properties: {
        name: { type: 'string', pattern: '[0-9a-zA-Z\-\_]+' },
        url: { type: 'string', format: 'uri' },
      },
    },
    default: [],
  },
  switchToNewStream: {
    type: 'boolean',
    default: true,
  },
  savedSessions: {
    type: 'number',
    default: 1,
  },
};

export default class CommonStore {
  @observable public streams = new Map<string, Stream>();
  @observable public apiPort = 0;
  @observable public apiVersion = '';

  @observable public switchToNewStream = true;
  @observable public savedSessions: string = '';

  private gqlClient?: GQLClient;
  private subscription?: ZenObservable.Subscription;

  @observable private streamSel?: string;
  @observable private navSel?: string;

  @observable private defaultPlugins = new Map<string, Plugin>();
  private diskStore = new Store({ schema });
  private autoSaveDispose?: IReactionDisposer;
  private timeUpdateTimeout?: NodeJS.Timeout;

  @action public init = async () => {
    const payload = await new Promise<RendererPayload>((resolve) => {
      ipcRenderer.once('renderer-payload', (_: IpcMessageEvent, arg: RendererPayload) => {
        resolve(arg);
      });
      ipcRenderer.send('renderer-payload');
    });

    this.apiPort = payload.apiPort;
    this.gqlClient = new GQLClient(payload.apiPort);
    if (payload.adminToken) {
      await this.gqlClient.initializeAdminToken(payload.adminToken, true);
    } else if (payload.adminOTP) {
      const adminToken = await this.gqlClient.initializeAdminToken(payload.adminOTP);
      ipcRenderer.send('save-admin-token', adminToken);
    }

    this.apiVersion = await this.gqlClient.getAPIVersion();

    const storeDefaultPlugins = this.diskStore.get('defaultPlugins') as Array<{name: string, url: string}> ;
    for (const { name, url } of storeDefaultPlugins) {
      await this.addPlugin('default', name, url);
    }

    this.switchToNewStream = this.diskStore.get('switchToNewStream') as boolean;
    this.savedSessions = String(this.diskStore.get('savedSessions'));

    this.autoSaveDispose = autorun(() => {
      const plugins = Array.from(this.defaultPlugins.values()).map((plugin) => ({
        name: plugin.name,
        url: plugin.url,
      }));
      this.diskStore.set('defaultPlugins', plugins);
      this.diskStore.set('switchToNewStream', this.switchToNewStream);
      this.diskStore.set('savedSessions', this.sanitizedSavedSessions);
    });

    this.timeUpdateTimeout = setInterval(() => {
      this.streams.forEach((stream) => stream.updateDisplayStartTime());
    }, 60000);

    this.subscription = await this.addStreamListeners(
      action(async (streamID: number) => {
        // tslint:disable-next-line: no-console
        console.log('Stream Created:', streamID);

        const newStream = new Stream(streamID);

        this.streams.set(newStream.uniqID, newStream);

        for (const [_, plugin] of this.defaultPlugins) {
          const pluginToken = await this.gqlAddPlugin(plugin.url);
          const newPlugin = plugin.withParams(newStream.uniqID, {
            apiURL: this.apiURL,
            apiToken: pluginToken,
            streamID,
          });
          newStream.plugins.set(newPlugin.id, newPlugin);
        }
        if (this.switchToNewStream) {
          this.selectStream(newStream.uniqID);
        }
      }),
      action(async (streamID: number) => {
        // tslint:disable-next-line: no-console
        console.log('Stream Removed', streamID);
        for (const [, stream] of this.streams) {
          if (stream.id === streamID) {
            stream.active = false;
            for (const [, plugin] of stream.plugins) {
              if (plugin.params) {
                await this.gqlRemovePlugin(plugin.params.apiToken);
              }
            }
          }
        }
        this.pruneInactiveStreams();
      }),
    );
  }

  @action public dispose = () => {
    if (this.autoSaveDispose) { this.autoSaveDispose(); }
    if (this.timeUpdateTimeout) { clearInterval(this.timeUpdateTimeout); }
    this.streams.forEach((stream) => {
      stream.plugins.forEach(async (plugin) => {
        if (plugin.params) {
          await this.gqlRemovePlugin(plugin.params.apiToken);
        }
      });
    });
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }

  @action public selectStream(streamID: string) {
    this.streamSel = streamID;
  }

  @action public selectNav(nav: string) {
    this.navSel = nav;
  }

  @action public async addPlugin(streamUniqID: string, name: string, url: string) {
    if (streamUniqID === 'default') {
      const newPlugin = new Plugin(name, url);
      this.defaultPlugins.set(newPlugin.id, newPlugin);
    } else if (this.streams.has(streamUniqID)) {
      const stream = this.streams.get(streamUniqID)!;
      const pluginToken = await this.gqlAddPlugin(url);
      const newPlugin = new Plugin(name, url, streamUniqID, {
        apiURL: this.apiURL,
        apiToken: pluginToken,
        streamID: stream.id,
      });
      stream.plugins.set(newPlugin.id, newPlugin);
    }
  }

  @action public removePlugins(plugins: Array<{ id: string, streamUniqID: string }>) {
    plugins.forEach(async ({ id, streamUniqID }) => {
      if (streamUniqID === 'default') {
        this.defaultPlugins.delete(id);
      } else if (this.streams.has(streamUniqID)) {
        const streamPlugins = this.streams.get(streamUniqID)!.plugins;
        if (streamPlugins.has(id)) {
          const plugin = streamPlugins.get(id)!;
          if (plugin.params) {
            await this.gqlRemovePlugin(plugin.params.apiToken);
          }
        }
        streamPlugins.delete(id);
      }
    });
  }

  @computed public get apiURL(): string {
    return `http://localhost:${this.apiPort}/query`;
  }
  @computed public get selectedStream(): Stream | undefined {
    if (this.streams.size > 0) {
      if (this.streamSel && this.streams.has(this.streamSel)) {
        return this.streams.get(this.streamSel)!;
      } else {
        return this.streams.values().next().value;
      }
    }
    return undefined;
  }

  @computed public get selectedStreamPlugins(): Plugin[] | undefined {
    if (!this.selectedStream) { return undefined; }
    return Array.from(this.selectedStream.plugins.values());
  }

  @computed public get selectedNav(): string {
    if (this.navSel === 'nav-settings') { return this.navSel; }
    if (this.selectedStream && this.selectedStream.plugins.size > 0) {
      if (!this.navSel || !this.navSel.startsWith(`plugin-${this.selectedStream.uniqID}`)) {
        return this.selectedStream.plugins.values().next().value.id;
      }
      return this.navSel;
    }
    return 'nav-settings';
  }

  @computed public get pluginsTree(): TreeViewSpec[] {
    const pluginsToTree = (streamUniqID: string, plugins: Map<string, Plugin>) =>
      Array.from(plugins).map(([_, plugin]) => ({
        id: plugin.id,
        name: `${plugin.name} - ${plugin.url}`,
        icon: plugin.icon,
        streamUniqID,
      }));

    const loadedPlugins = Array.from(this.streams).map(([streamUniqID, stream]) => ({
      id: streamUniqID,
      name: stream.displayName,
      children: pluginsToTree(streamUniqID, stream.plugins),
    }));

    return loadedPlugins.concat([{
      id: 'default',
      name: `Default Plugins`,
      children: pluginsToTree('default', this.defaultPlugins),
    }]);
  }

  @computed public get sanitizedSavedSessions(): number {
    const savedSessions = Number(this.savedSessions);
    if (!this.savedSessions || isNaN(savedSessions)) {
      return 1;
    }
    return savedSessions;
  }

  @computed public get displaySavedSessions(): string {
    if (this.sanitizedSavedSessions < 0) {
      return 'infinite';
    }
    return `${this.sanitizedSavedSessions}`;
  }

  @action private pruneInactiveStreams() {
    const candidatesForRemoval: string[] = [];
    for (const [streamUniqID, stream] of this.streams) {
      if (!stream.active) {
        candidatesForRemoval.push(streamUniqID);
      }
    }
    const sss = this.sanitizedSavedSessions;
    if (sss >= 0 && sss < candidatesForRemoval.length) {
      const toRemove = (sss === 0) ? candidatesForRemoval : candidatesForRemoval.slice(0, -sss);
      for (const streamUniqID of toRemove) {
        this.streams.delete(streamUniqID);
      }
    }
  }

  private gqlAddPlugin = (pluginURL: string): Promise<string> => {
    if (!this.gqlClient) { throw new Error('Must start GQL client first.'); }
    return this.gqlClient.addPlugin(pluginURL);

  }

  private gqlRemovePlugin = (apiToken: string): Promise<boolean> => {
    if (!this.gqlClient) { throw new Error('Must start GQL client first.'); }
    return this.gqlClient.removePlugin(apiToken);
  }

  private addStreamListeners = async (
    onStreamCreated: (streamID: number) => void,
    onStreamRemoved: (streamID: number) => void,
  ) => {
    if (!this.gqlClient) { throw new Error('Must start GQL client first.'); }
    const streams = await this.gqlClient.listStreams();
    streams.map(({ id }) => onStreamCreated(id));
    return this.gqlClient.subscribeToStreamEvents((streamID, type) => {
      switch (type) {
        case 'AddStream':
          onStreamCreated(streamID);
          break;
        case 'RemoveStream':
          onStreamRemoved(streamID);
          break;
      }
    });
  }
}
