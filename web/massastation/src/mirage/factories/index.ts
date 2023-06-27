import { accountFactory } from './account';
import { domainFactory } from './domain';
import { networkFactory } from './network';
import { pluginFactory } from './plugin';
import { storeFactory } from './store';

export const factories = {
  account: accountFactory,
  domain: domainFactory,
  network: networkFactory,
  plugin: pluginFactory,
  store: storeFactory,
};
