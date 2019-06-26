<template>
  <v-app dark>
    <v-app-nav-drawer
      width=260
      height="calc(100vh - 30px)"
      :state="state"
    />
    <v-content>
      <v-window v-model="state.selectedNav" vertical>
        <template v-for="stream in state.streams.values()">
          <v-window-item
            v-for="[_, plugin] in stream.plugins"
            :key="plugin.id"
            :value="plugin.id"
            :transition="false"
            :reverseTransition="false"
          >
            <div class="tab-item-wrapper">
              <WebView :plugin="plugin" />
            </div>
          </v-window-item>
        </template>
        <v-window-item
          key="settings"
          value="nav-settings"
          :transition="false"
          :reverseTransition="false"
        >
          <div class="tab-item-wrapper">
            <Settings :state="state" />
          </div>
        </v-window-item>
      </v-window>
    </v-content>
  </v-app>
</template>

<script>
import { remote, ipcRenderer } from 'electron';
import WebView from './WebView';
import NavDrawer from './NavDrawer';
import Settings from './Settings';
import { Titlebar, Color } from 'custom-electron-titlebar';

import { observer } from 'mobx-vue';
import CommonStore from '../stores/commonStore';

export default observer({
  name: 'App',
  // Locally register these components
  components: {
    WebView,
    'v-app-nav-drawer': NavDrawer,
    Settings,
  },
  data: () => {
    return {
      state: new CommonStore(),
    };
  },
  mounted() {
    this.titlebar = new Titlebar({
      backgroundColor: Color.fromHex('#212121'),
      icon: 'icon.ico',
      menu: new remote.Menu(),
    });
    this.titlebar.updateTitle('Aetherometer');
    this.state.init();
    window.addEventListener('beforeunload', () => {
      this.titlebar.dispose();
      this.state.dispose();
      // Intentionally block the renderer process so that we have time to
      // clean up.
      ipcRenderer.sendSync('unloading');
    });
  },
  beforeDestroy() {
    this.titlebar.dispose();
    this.state.dispose();
    // Intentionally block the renderer process so that we have time to
    // clean up.
    ipcRenderer.sendSync('unloading');
  },
});
</script>

<style>
  html { overflow-y: auto }
  .application--wrap {
    min-height: calc(100vh - 30px);
  }
  .titlebar {
    font-family: 'Roboto', sans-serif;
    font-weight: 400;
  }
  .tab-item-wrapper {
    height: calc(100vh - 30px);
    display: flex;
  }
</style>