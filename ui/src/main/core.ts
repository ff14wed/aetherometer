import { app } from 'electron';
import fs from 'fs';
import path from 'path';
import os from 'os';
import ChildProcess from 'child_process';

import appRootDir from 'app-root-dir';
import * as upath from 'upath';

import cryptoRandomString from 'crypto-random-string';
import getPort from 'get-port';

const getOS = () => {
  switch (process.platform) {
    case 'darwin':
      return 'mac';
    case 'win32':
      return 'win';
    default:
      return 'linux';
  }
};

const isDevelopment = process.env.NODE_ENV !== 'production';

export default class Core {
  private userData = app.getPath('userData');
  private logsPath = app.getPath('logs');
  private configPath = path.join(this.userData, 'core-config.toml');
  private commandOutLogs = path.join(this.logsPath, 'core.out.log');
  private commandErrLogs = path.join(this.logsPath, 'core.err.log');
  private datasheetsPath = upath.toUnix(path.join(this.resourcesPath, 'datasheets'));
  private mapCachePath = upath.toUnix(path.join(this.resourcesPath, 'maps'));

  private corePort = 0;
  private adminOTP = cryptoRandomString({ length: 16, type: 'base64' });
  private adminToken?: string;


  private get resourcesPath() {
    return isDevelopment ?
      path.join(appRootDir.get(), '..', 'resources') :
      path.join(appRootDir.get());
  }

  private get binDirPath() {
    return isDevelopment ?
      path.join(appRootDir.get(), '..', 'resources', getOS()) :
      path.join(appRootDir.get(), 'bin');
  }

  private get coreBinPath() {
    const coreBin = (process.platform === 'win32') ? 'core.exe' : 'core';
    return path.join(this.binDirPath, coreBin);
  }

  private get pkillPath() {
    return path.join(this.binDirPath, 'windows-kill.exe');
  }

  private get hookDLLPath() {
    return upath.toUnix(path.join(this.binDirPath, 'xivhook.dll'));
  }

  public start = async () => {
    this.corePort = await getPort({ host: '127.0.0.1', port: [8080, 8081, 8082] });
    await this.writeConfigFile();

    const cmdOutStream = fs.createWriteStream(this.commandOutLogs, { flags: 'a+' });
    const cmdErrStream = fs.createWriteStream(this.commandErrLogs, { flags: 'a+' });

    const child = ChildProcess.spawn(this.coreBinPath, ['-c', this.configPath]);
    child.stdout.setEncoding('utf8');

    const closeCallback = (code: number | null, signal: string | null) => {
      cmdErrStream.write(`Command exited with code ${code} or signal ${signal}${os.EOL}`);
    };

    child.stderr.pipe(cmdErrStream, { end: false });

    await new Promise((resolve, reject) => {
      child.stdout.on('data', (data) => {
        cmdOutStream.write(data.toString());
        if (/http-server.*Running/.test(data.toString())) {
          resolve();
        }
      });

      child.on('close', (code, signal) => {
        closeCallback(code, signal);
        reject(`Core exited with code ${code}. Check logs at ${this.commandOutLogs} for details.`);
      });
    });


    child.stdout.removeAllListeners('data');
    child.stdout.pipe(cmdOutStream, { end: false });

    child.stderr.removeAllListeners('close');
    child.on('close', closeCallback);


    this.exitExecutor = (exitResolve) => {
      child.removeAllListeners('close');
      const conds = [
        new Promise((resolve) => {
          child.on('exit', async (code, signal) => {
            closeCallback(code, signal);
            resolve();
          });
        }),
        new Promise((resolve) => {
          child.stdout.on('end', () => {
            resolve();
          });
        }),
        new Promise((resolve) => {
          child.stderr.on('end', () => {
            resolve();
          });
        }),
      ];

      if (process.platform === 'win32') {
        ChildProcess.execSync(`${this.pkillPath} -SIGINT ${child.pid}`);
      } else {
        child.kill('SIGINT');
      }
      Promise.all(conds).then(() => {
        cmdOutStream.end();
        cmdErrStream.end();
        exitResolve();
      });
    };
  }

  public getRendererPayload = () => {
    if (this.adminToken) {
      return { apiPort: this.corePort, adminToken: this.adminToken };
    }
    return { apiPort: this.corePort, adminOTP: this.adminOTP };
  }

  public saveAdminToken = (token: string) => {
    this.adminToken = token;
  }

  public exit = async () => {
    await new Promise(this.exitExecutor);
  }

  private exitExecutor = (resolve: () => void) => resolve();

  private writeConfigFile = async () => {
    const mapDirPromise = new Promise((resolve, reject) => {
      fs.access(this.mapCachePath, fs.constants.F_OK, (err) => {
        if (err) {
          fs.mkdir(this.mapCachePath, (mkdirErr) => {
            if (mkdirErr) { reject(mkdirErr); } else { resolve(); }
          });
        } else {
          resolve();
        }
      });
    });

    const configStream = fs.createWriteStream(this.configPath);
    configStream.write(`api_port = ${this.corePort}` + os.EOL);
    configStream.write(`data_path = "${this.datasheetsPath}"` + os.EOL);
    configStream.write(`admin_otp = "${this.adminOTP}"` + os.EOL);
    configStream.write('[maps]' + os.EOL);
    configStream.write(`cache = "${this.mapCachePath}"` + os.EOL);
    configStream.write('[adapters.hook]' + os.EOL);
    configStream.write('enabled = true' + os.EOL);
    configStream.write(`dll_path = "${this.hookDLLPath}"` + os.EOL);
    configStream.write('ffxiv_process = "ffxiv_dx11.exe"' + os.EOL);

    const configPromise = new Promise((resolve) => {
      configStream.end(resolve);
    });
    await Promise.all([mapDirPromise, configPromise]);
  }
}
