<template>
  <v-app>
    <v-main>
      <v-container fluid fill-height>
        <v-layout align-center column justify-center>
          <h1 class="display-2 font-weight-thin mb-3">Aetherometer</h1>
          <h4 class="subheading statusline">{{ status }}</h4>
          <v-btn color="error" v-if="displayClose" v-on:click="closeApp">Close</v-btn>
        </v-layout>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import 'typeface-roboto';

import { ipcRenderer } from 'electron';

export default {
  name: 'Loading',
  data: () => {
    return {
      status: '',
      displayClose: false,
    };
  },
  methods: {
    closeApp: (event) => ipcRenderer.send('close-from-loading'),
  },
  mounted() {
    this.$nextTick(() => {
      ipcRenderer.on('status', (e, arg, displayClose) => {
        this.status = arg;
        this.displayClose = displayClose;
      });
    });
  },
};
</script>

<style>
html {
  overflow-y: hidden;
}
.statusline {
  word-break: break-all;
}
</style>
