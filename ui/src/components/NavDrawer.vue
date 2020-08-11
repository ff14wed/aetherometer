<template>
  <v-navigation-drawer
    id="nav-drawer"
    class="elevation-10 grey darken-4"
    style="margin-top: 30px"
    app
    floating
    permanent
    :mini-variant="mini"
    mini-variant-width="60"
    width="260"
    height="calc(100vh - 30px)"
  >
    <v-layout column fill-height>
      <div id="plugin-list">
        <v-layout tag="v-list" column>
          <v-nav-list-item
            v-for="plugin in state.selectedStreamPlugins"
            :key="plugin.id"
            :navID="plugin.id"
            :name="plugin.name"
            :showTooltip="mini"
            :state="state"
          >
            <img
              v-if="plugin.icon"
              :src="plugin.icon"
              :alt="plugin.id"
              @error="onFaviconErr(plugin)"
            />
            <v-icon v-else large>extension</v-icon>
          </v-nav-list-item>
        </v-layout>
      </div>

      <v-spacer />

      <v-divider />
      <v-list>
        <v-nav-list-item navID="nav-settings" name="Settings" :showTooltip="mini" :state="state">
          <v-icon large>settings</v-icon>
        </v-nav-list-item>

        <v-menu offset-x right v-if="state.selectedStream">
          <template v-slot:activator="{ on }">
            <v-list-item @click="1" v-on="on">
              <v-nav-icon showTooltip :tooltip="state.selectedStream.displayName">
                <v-icon large>input</v-icon>
              </v-nav-icon>
              <v-list-item-title class="font-weight-light">{{ state.selectedStream.shortName }}</v-list-item-title>
            </v-list-item>
          </template>
          <v-list>
            <v-list-item
              v-for="[streamUniqID, stream] in state.streams"
              :key="streamUniqID"
              @click="state.selectStream(streamUniqID)"
            >
              <v-list-item-title>{{ stream.displayName }}</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>

        <v-list-item>
          <v-list-item-action>
            <v-btn icon @click.stop="mini = !mini">
              <v-icon large v-if="mini">chevron_right</v-icon>
              <v-icon large v-else>chevron_left</v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </v-list>
    </v-layout>
  </v-navigation-drawer>
</template>

<script>
import { Titlebar, Color } from 'custom-electron-titlebar';

import { action, observer } from 'mobx-vue';
import CommonStore from '../stores/commonStore';
import NavListItem from './NavListItem.vue';
import NavIcon from './NavIcon.vue';

export default observer({
  props: {
    state: {
      type: CommonStore,
      required: true,
    },
    width: {
      type: [Number, String],
      default: 300,
    },
    height: {
      type: [Number, String],
      default: '100%',
    },
  },
  components: {
    'v-nav-list-item': NavListItem,
    'v-nav-icon': NavIcon,
  },
  data: () => ({
    mini: false,
  }),
  methods: {
    onFaviconErr: (plugin) => {
      plugin.setIcon(null);
    },
  },
});
</script>

<style>
#plugin-list {
  max-height: 80vh;
  overflow-y: auto;
  overflow-x: hidden;
  color: rgba(255, 255, 255, 0);
  transition: color 0.3s ease;
}

#plugin-list:hover {
  color: rgba(255, 255, 255, 0.2);
}

#plugin-list::-webkit-scrollbar {
  width: 5px;
}

#plugin-list::-webkit-scrollbar-thumb {
  /* box-shadow: inset 0 0 0 10px; */
  box-shadow: inset 0 0 0 10px;
}
</style>
