'use strict';

import { app, protocol, BrowserWindow, ipcMain, IpcMainEvent } from 'electron';
import {
  createProtocol,
  installVueDevtools,
} from 'vue-cli-plugin-electron-builder/lib';

import Store from 'electron-store';

import Core from './main/core';

const isDevelopment = process.env.NODE_ENV !== 'production';

// Keep a global reference of the window object, if you don't, the window will
// be closed automatically when the JavaScript object is garbage collected.
let win: BrowserWindow | null;

interface StoreType {
  width: number;
  height: number;
}

const store = new Store<StoreType>({
  defaults: {
    width: 1920,
    height: 1080,
  },
});

const coreInst = new Core();

const quitApp = async () => {
  // tslint:disable-next-line: no-console
  console.log('waiting for core to exit');
  await coreInst.exit();
  // tslint:disable-next-line: no-console
  console.log('done waiting for core');
  app.quit();
};

let createdAppProtocol = false;
const createAppProtocol = () => {
  if (!createdAppProtocol) {
    createProtocol('app');
    createdAppProtocol = true;
  }
};

// Scheme must be registered before the app is ready
protocol.registerSchemesAsPrivileged([{ scheme: 'app', privileges: { secure: true, standard: true } }]);

const createLoadingWindow = async () => {
  const loadingWin = new BrowserWindow({
    width: 400,
    height: 300,
    show: false,
    frame: false,
    resizable: false,
    webPreferences: {
      nodeIntegration: true,
      devTools: true,
    },
  });

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    // Load the url of the dev server if in development mode
    loadingWin.loadURL(`${process.env.WEBPACK_DEV_SERVER_URL as string}/loading.html`);
  } else {
    createAppProtocol();
    loadingWin.loadURL('app://./loading.html');
  }

  ipcMain.once('close-from-loading', () => {
    loadingWin.close();
  });

  return new Promise<BrowserWindow>((resolve) => {
    loadingWin.once('ready-to-show', () => {
      loadingWin.show();
      resolve(loadingWin);
    });
  });
};

const createMainWindow = async () => {
  // Create the browser window.
  win = new BrowserWindow({
    width: store.get('width'),
    height: store.get('height'),
    x: store.get('x'),
    y: store.get('y'),
    show: false,
    frame: false,
    webPreferences: {
      nodeIntegration: true,
      webviewTag: true,
      devTools: true,
    },
  });

  const window = win;

  ipcMain.on('save-admin-token', (event: IpcMainEvent, token: string) => {
    coreInst.saveAdminToken(token);
  });

  ipcMain.on('renderer-payload', (event: IpcMainEvent) => {
    event.sender.send('renderer-payload', coreInst.getRendererPayload());
  });

  ipcMain.on('unloading', (event: IpcMainEvent) => {
    setTimeout(() => {
      event.returnValue = true;
    }, 100);
  });

  if (process.env.WEBPACK_DEV_SERVER_URL) {
    // Load the url of the dev server if in development mode
    win.loadURL(`${process.env.WEBPACK_DEV_SERVER_URL as string}/index.html`);
    if (!process.env.IS_TEST) { win.webContents.openDevTools({ mode: 'right' }); }
  } else {
    createAppProtocol();
    // Load the index.html when not in development
    win.loadURL('app://./index.html');
  }


  win.on('close', () => {
    const windowState = window.getBounds();
    store.set('width', windowState.width);
    store.set('height', windowState.height);
    store.set('x', windowState.x);
    store.set('y', windowState.y);
  });

  win.on('closed', () => {
    win = null;
  });

  return new Promise((resolve) => {
    window.once('ready-to-show', () => {
      window.show();
      resolve();
    });
  });
};

// Quit when all windows are closed.
app.on('window-all-closed', async () => {
  // On macOS it is common for applications and their menu bar
  // to stay active until the user quits explicitly with Cmd + Q
  if (process.platform !== 'darwin') {
    await quitApp();
  }
});

app.on('activate', () => {
  // On macOS it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (win === null) {
    createMainWindow();
  }
});

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', async () => {
  if (isDevelopment && !process.env.IS_TEST) {
    // Install Vue Devtools
    try {
      await installVueDevtools();
    } catch (e) {
      // tslint:disable-next-line: no-console
      console.error('Vue Devtools failed to install:', e.toString());
    }
  }
  const loadingWin = await createLoadingWindow();
  try {
    loadingWin.webContents.send('status', 'Starting core...');
    await coreInst.start();

    loadingWin.webContents.send('status', 'Loading UI...');
    await createMainWindow();

    loadingWin.close();
    ipcMain.removeAllListeners('close-from-loading');


  } catch (err) {
    loadingWin.webContents.send('status', `Error encountered: ${err}`, true);
  }
});

// Exit cleanly on request from parent process in development mode.
if (isDevelopment) {
  if (process.platform === 'win32') {
    process.on('message', async (data) => {
      if (data === 'graceful-exit') {
        await quitApp();
      }
    });
  } else {
    process.on('SIGTERM', async () => {
      await quitApp();
    });
  }
}
