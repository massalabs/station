import { routesForAccounts } from './account';
import { routesForLogs } from './logs';
import { routesForNetwork } from './network';
import { routesForPlugins } from './plugin';
import { routesForStore } from './store';
import { routesForVersion } from './version';

const handlers = {
  accounts: routesForAccounts,
  logs: routesForLogs,
  network: routesForNetwork,
  plugins: routesForPlugins,
  store: routesForStore,
  version: routesForVersion,
};

export { handlers };
