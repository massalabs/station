import { createServer } from 'miragejs';
import { handlers } from './handlers';
import { models } from './models';
import { factories } from './factories';
import { ENV } from '../const/env/env';

export function mockServer(environment = ENV.DEV) {
  const server = createServer({
    environment,
    models,
    factories,
    seeds(server) {
      server.createList('plugin', 2);
      server.createList('domain', 50);
      server.createList('store', 7);
      server.createList('account', 5);
      server.create('network');
    },
  });

  for (const namespace of Object.keys(handlers)) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    handlers[namespace](server);
  }

  return server;
}
