import { Server } from 'miragejs';
import { AppSchema } from '../types';

export function routesForDomains(server: Server) {
  server.get(
    '/all/domains',
    (schema: AppSchema) => {
      let { models: domains } = schema.all('domain');

      return domains;
    },
    { timing: 5000 },
  );
}
