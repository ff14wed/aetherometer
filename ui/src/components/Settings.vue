<template>
  <div id="settings-container">
    <v-container fill-height>
      <v-layout column>
        <h5 class="section-heading headline font-weight-thin">About</h5>
          <v-sheet class="mb-3" elevation=16>
            <v-container>
              <v-list>
                <v-list-tile>
                  <v-list-tile-content>Aetherometer Version:</v-list-tile-content>
                  <v-list-tile-content class="align-end">v{{ getAppVersion() }}</v-list-tile-content>
                </v-list-tile>
                <v-list-tile>
                  <v-list-tile-content>Aetherometer API version:</v-list-tile-content>
                  <v-list-tile-content class="align-end">{{ state.apiVersion }}</v-list-tile-content>
                </v-list-tile>
                <v-list-tile>
                  <v-list-tile-content>Aetherometer API URL:</v-list-tile-content>
                  <v-list-tile-content class="align-end">{{ state.apiURL }}</v-list-tile-content>
                </v-list-tile>
                <v-list-tile>
                  <v-list-tile-content>Debug Logs:</v-list-tile-content>
                  <v-list-tile-content class="align-end">
                    <v-btn @click="goToLogsPath()">View Core Logs</v-btn>
                  </v-list-tile-content>
                </v-list-tile>
                <v-list-tile>
                  <v-list-tile-content>Help / Docs / Source Code:</v-list-tile-content>
                  <v-list-tile-content class="align-end">
                    <v-layout>
                      <v-flex><v-btn @click="goToLink('https://github.com/ff14wed/aetherometer')">GitHub</v-btn></v-flex>
                    </v-layout>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
            </v-container>
          </v-sheet>

        <h5 class="section-heading headline font-weight-thin">Manage Plugins</h5>
        <v-sheet class="mb-3" elevation=16>
          <v-container>
            <v-layout column>
              <v-treeview
                v-model="tree"
                :items="state.pluginsTree"
                open-all
                hoverable
                selectable
                return-object
                open-on-click
              >
                <template v-slot:prepend="{ item }">
                  <img v-if="item.icon" :src="item.icon" width=24 height=24 />
                  <v-icon v-else-if="item.streamUniqID">extension</v-icon>
                </template>
              </v-treeview>
              <v-flex align-self-end class="pb-3">
                <v-btn @click="unselectAll">Unselect All</v-btn>
                <v-btn :disabled="selectedPlugins.length === 0" @click="removeSelectedPlugins">
                  <template v-if="selectedPlugins.length === 1">Remove Plugin</template>
                  <template v-else>Remove Plugins</template>
                </v-btn>
              </v-flex>

              <v-divider class="pb-3" />

              <v-form ref="addPluginForm" lazy-validation v-model="pluginFormValid">
                <v-subheader class="pa-0">Add Plugin to Selected Streams</v-subheader>
                <v-layout row justify-space-between>
                  <v-flex xs12 md4>
                    <v-text-field
                      v-model="addPluginName"
                      :rules="pluginNameRules"
                      label="Plugin Name"
                      required
                    />
                  </v-flex>
                  <v-flex xs12 md7>
                    <v-text-field
                      v-model="addPluginURL"
                      :rules="pluginURLRules"
                      label="Plugin URL"
                      required
                    />
                  </v-flex>
                </v-layout>
                <v-layout row justify-end>
                  <v-input readonly :error-messages="pluginStreamErrors">
                    {{ selectedStreams.length }} Streams Selected
                  </v-input>
                  <v-btn
                    :disabled="!pluginFormValid"
                    color="success"
                    @click="formAddPlugin"
                  >
                    Add Plugin
                  </v-btn>

                  <v-btn
                    color="error"
                    @click="formReset"
                  >
                    Reset Form
                  </v-btn>
                </v-layout>
              </v-form>
            </v-layout>
          </v-container>
        </v-sheet>

        <h5 class="section-heading headline font-weight-thin">Miscellaneous</h5>
          <v-sheet class="mb-3" elevation=16>
            <v-container>
              <v-layout column>
                <v-switch
                  v-model="state.switchToNewStream"
                  :label="`Automatically switch to new session when a stream is created`"
                ></v-switch>
                <div class="pb-2 font-weight-light red--text text--lighten-1">
                  Warning: Be sure to save any important data before allowing stream sessions
                  to be closed automatically.
                </div>
                <v-text-field
                  v-model="state.savedSessions"
                  :label="`${state.displaySavedSessions} inactive stream session(s) retained. Specify a negative number to retain all inactive stream sessions.`"
                  type="number"
                ></v-text-field>
              </v-layout>
            </v-container>
          </v-sheet>
      </v-layout>
    </v-container>
  </div>
</template>

<script>
import { remote, shell } from 'electron';

import CommonStore from '../stores/commonStore';

import { observer } from 'mobx-vue';

import validURL from 'valid-url';

const nameRegex = /^[0-9a-zA-Z\-\_]+$/;

const isNameValid = (v) => {
  if (!v) { return false; }
  return !!v.match(nameRegex) || false;
};

export default observer({
  props: {
    state: {
      type: CommonStore,
      required: true,
    },
  },
  data: () => ({
    tree: [],
    pluginFormValid: true,
    addPluginName: '',
    addPluginURL: '',
    pluginNameRules: [
      (v) => !!v || 'Plugin Name is required',
      (v) => isNameValid(v) || 'Plugin Name must be alphanumeric with no spaces (hyphens and underscore allowed)',
    ],
    pluginURLRules: [
      (v) => !!v || 'Plugin URL is required',
      (v) => !!validURL.isWebUri(v) || 'Plugin URL must be valid (must include http or https scheme)',
    ],
  }),
  computed: {
    // Basically plugin items have the streamUniqID property, and streams don't
    selectedPlugins() {
      return this.tree.filter((item) => item.streamUniqID);
    },
    selectedStreams() {
      return this.tree.filter((item) => !item.streamUniqID);
    },
    pluginStreamErrors() {
      if ((this.addPluginName && this.addPluginName.length > 0) ||
          (this.addPluginURL && this.addPluginURL.length > 0)) {
        if (this.selectedStreams.length === 0) {
          return ['At least one selected stream is required to add plugin'];
        }
      }
      return [];
    },
  },
  methods: {
    removeSelectedPlugins() {
      this.state.removePlugins(this.selectedPlugins);
    },
    formAddPlugin() {
      if (this.$refs.addPluginForm.validate()) {
        for (const { id } of this.selectedStreams) {
          this.state.addPlugin(id, this.addPluginName, this.addPluginURL);
        }
        this.$refs.addPluginForm.reset();
      }
    },
    formReset() {
      this.$refs.addPluginForm.reset();
    },
    getAppVersion() {
      return remote.app.getVersion();
    },
    unselectAll() {
      this.tree = [];
    },
    goToLogsPath() {
      shell.openItem(remote.app.getPath('logs'));
    },
    goToLink(url) {
      shell.openExternal(url);
    },
  },
});
</script>

<style scoped>
  .section-heading {
    margin-bottom: 10px;
  }

  #settings-container {
    overflow-y: auto;
    width: 100%;
    height: 100%;
    flex: 1;
  }
  #settings-container::-webkit-scrollbar {
    width: 5px;
  }

  #settings-container::-webkit-scrollbar-thumb {
    background-color: rgba(255, 255, 255, 0.2);
  }
</style>
