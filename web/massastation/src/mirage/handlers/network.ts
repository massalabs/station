import { Server, Response } from 'miragejs';
import { AppSchema } from '../types';

export function routesForNetwork(server: Server) {
  server.get('/network', (schema: AppSchema) => {
    const { models: network } = schema.all('network');

    return network;
  });

  server.post('/network/:network', (schema: AppSchema, request) => {
    const { network } = request.params;
    const storedNetwork = schema.find('network', network);

    if (!storedNetwork)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    storedNetwork.update({ currentNetwork: request.params.network });

    return new Response(200, {});
  });
}
