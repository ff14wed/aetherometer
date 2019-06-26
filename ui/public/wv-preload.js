var webFrame = require('electron').webFrame;

webFrame.executeJavaScript('window.waitForInit = true');
