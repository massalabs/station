import { Server } from 'miragejs';
import { AppSchema } from '../types';

export function routesForAccounts(server: Server) {
  server.get(
    '/plugin/massa-labs/massa-wallet/api/accounts',
    (schema: AppSchema) => {
      let { models: accounts } = schema.all('account');

      return accounts;
    },
    { timing: 500 },
  );
}
