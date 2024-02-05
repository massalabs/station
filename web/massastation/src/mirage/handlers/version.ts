import { Server } from 'miragejs';

export function routesForVersion(server: Server) {
  server.get('/version', () => {
    return 'dev';
  });
}
