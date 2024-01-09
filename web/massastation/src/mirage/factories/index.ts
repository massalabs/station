import { accountFactory } from './account';
import { networkFactory } from './network';
import { pluginFactory } from './plugin';
import { storeFactory } from './store';

export const factories = {
  account: accountFactory,
  network: networkFactory,
  plugin: pluginFactory,
  store: storeFactory,
};
