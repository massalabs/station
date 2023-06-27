import { routesForAccounts } from './account';
import { routesForDomains } from './domain';
import { routesForNetwork } from './network';
import { routesForPlugins } from './plugin';
import { routesForStore } from './store';

const handlers = {
  accounts: routesForAccounts,
  domains: routesForDomains,
  network: routesForNetwork,
  plugins: routesForPlugins,
  store: routesForStore,
};

export { handlers };
