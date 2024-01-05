import { routesForAccounts } from './account';
import { routesForNetwork } from './network';
import { routesForPlugins } from './plugin';
import { routesForStore } from './store';

const handlers = {
  accounts: routesForAccounts,
  network: routesForNetwork,
  plugins: routesForPlugins,
  store: routesForStore,
};

export { handlers };
