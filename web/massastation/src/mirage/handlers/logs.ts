import { Response, Server } from 'miragejs';

const levels = ['DEBUG', 'INFO', 'INFO', 'INFO', 'WARN', 'ERROR'];

export function routesForLogs(server: Server) {
  server.get('/logs', (_, request) => {
    const lines = Number(request.queryParams.lines) || 1500;
    const start = Date.now() - lines * 1000;

    const logs = Array.from({ length: lines }, (_, i) => {
      const date = new Date(start + i * 1000).toISOString();
      const level = levels[i % levels.length];
      return `${date}\t${level}\tdev log line ${i}`;
    }).join('\n');

    return new Response(200, { 'Content-Type': 'text/plain' }, `${logs}\n`);
  });
}
