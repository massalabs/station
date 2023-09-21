export interface MassaStoreModel {
  name: string;
  author: string;
  description: string;
  version: string;
  url: string;
  logo: string;
  file: pluginStoreItemFile;
  os: string;
  iscompatible: boolean;
  massastationMinVersion: string;
}

export interface pluginStoreItemFile {
  url: string;
  checksum: string;
}

export enum PluginStatus {
  Down = 'Down',
  Up = 'Up',
  Starting = 'Starting',
  Stopping = 'Stopping',
  Error = 'Error',
  NotInstalled = 'NotInstalled',
  Crashed = 'Crashed',
}
