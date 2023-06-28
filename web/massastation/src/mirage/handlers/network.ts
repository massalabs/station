import { Server, Response } from 'miragejs';
import { AppSchema } from '../types';

export function routesForNetwork(server: Server) {
  server.get('/network', (schema: AppSchema) => {
    const network = schema.findBy('network', { id: '1' });

    if (!network)
      return new Response(400, {}, { code: '400', error: 'Bad request' });

    return network.attrs;
  });

  server.post('/network/:network', (schema: AppSchema, request) => {
    const { network } = request.params;
    const storedNetwork = schema.findBy('network', { id: '1' });

    if (!storedNetwork)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    storedNetwork.update({ currentNetwork: network });

    return new Response(200, {});
  });
}
