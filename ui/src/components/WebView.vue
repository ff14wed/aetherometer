<template>
  <webview
    v-bind:id="plugin.id"
    v-bind:src="plugin.url"
    v-bind:style="styleObject"
    v-bind:preload="preload"
    ref="webview"
    webpreferences="contextIsolation"
  ></webview>
</template>

<script>
import path from 'path';
import url from 'url';
import { shell, remote, webContents } from 'electron';

import Plugin from '../stores/plugin';

const isInternalURL = (cur, target) => {
  const curURL = url.parse(cur);
  const testURL = url.parse(target);
  return (curURL.protocol === testURL.protocol) &&
    (curURL.auth === testURL.auth) &&
    (curURL.host === testURL.host) &&
    (curURL.path === testURL.path);
};

const isWebURL = (destURL) => {
  return destURL.startsWith('http:') || destURL.startsWith('https:');
};

const navigateToExternalURL = (src, destURL) => {
  if (isWebURL(destURL)) {
    const dialogOptions = {
      buttons: ['Cancel', 'Continue'],
      title: 'External URL Confirmation',
      message: 'Go to external URL?',
      detail: `Page at ${src} wants to navigate to ${destURL}. Proceed?`,
      type: 'question',
    };

    remote.dialog.showMessageBox(dialogOptions, (response) => {
      if (response > 0) {
        shell.openExternal(url);
      }
    });
  }
};

export default {
  props: {
    plugin: {
      type: Plugin,
      required: true,
    },
  },
  data: () => {
    return {
      styleObject: {
        flexGrow: 1,
      },
      preload: url.pathToFileURL(path.join(__static, './wv-preload.js')),
    };
  },

  mounted() {
    const webview = this.$refs.webview;

    webview.addEventListener('page-favicon-updated', (event) => {
      if (event.favicons.length > 0) {
        this.plugin.setIcon(event.favicons[0]);
      }
    });

    webview.addEventListener('dom-ready', () => {
      webview.getWebContents().on('before-input-event', (event, input) => {
        if (input.type !== 'keyDown') {
          return;
        }

        if (input.control && input.shift && input.key === 'J') {
          webview.getWebContents().openDevTools({ mode: 'undocked' });
        }
      });

      webview.getWebContents().on('will-navigate', (event, destURL) => {
        if (isInternalURL(this.plugin.url, destURL)) {
          return;
        }
        event.preventDefault();
        webview.stop();
        webview.getWebContents().stop();
        navigateToExternalURL(this.plugin.url, destURL);
      });

      webview.getWebContents().executeJavaScript( `
        if (window.waitForInit && window.initPlugin) {
          window.initPlugin(${JSON.stringify(this.plugin.params)});
        }
      `);
    });
  },
};

</script>