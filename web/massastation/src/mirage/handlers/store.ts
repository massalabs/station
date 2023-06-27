import { Server } from 'miragejs';
import { AppSchema } from '../types';

export function routesForStore(server: Server) {
  server.get('plugin-store', (schema: AppSchema) => {
    let { models: stores } = schema.all('store');

    return stores;
  });
}
