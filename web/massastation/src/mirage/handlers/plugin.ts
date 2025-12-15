import { Server, Response } from 'miragejs';
import { AppSchema } from '../types';
import { faker } from '@faker-js/faker';

export function routesForPlugins(server: Server) {
  server.get('/plugin-manager', (schema: AppSchema) => {
    let { models: plugins } = schema.all('plugin');

    return plugins;
  });

  server.get('/plugin-manager/:id', (schema: AppSchema, request) => {
    const { id } = request.params;

    let plugin = schema.find('plugin', id);

    if (!plugin)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    return plugin.attrs;
  });

  server.post('/plugin-manager/:id/execute', (schema: AppSchema, request) => {
    const { command } = JSON.parse(request.requestBody);
    const { id } = request.params;
    const plugin = schema.find('plugin', id);

    if (!plugin)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    const status = ['update', 'start'].includes(command) ? 'Up' : 'Down';
    const updatable = command === 'update' ? false : plugin.updatable;

    plugin.update({
      version: faker.system.semver(),
      status,
      updatable,
    });

    return new Response(200, {});
  });

  server.post('plugin-manager', (schema: AppSchema, request) => {
    const sourceParam = request.queryParams.source;
    const sourceURL = Array.isArray(sourceParam) ? sourceParam[0] : sourceParam;

    if (!sourceURL) {
      return new Response(400, {}, { code: '400', error: 'Missing source URL' });
    }

    const sourcePath = sourceURL.split('?')[0];

    const storePlugin = schema.findBy(
      'store',
      (store) => store.file.url === sourcePath,
    );
    const newPlugin = {
      id: faker.number.int().toString(),
      status: 'Up',
      updatable: false,
      name: storePlugin?.name || faker.lorem.word(),
      author: storePlugin?.author || faker.person.firstName(),
      description: storePlugin?.description || faker.lorem.sentence(),
      version: storePlugin?.version || faker.system.semver(),
      home: storePlugin?.url || `/plugin/massalabs/${faker.lorem.word()}/`,
    };
    schema.create('plugin', newPlugin);

    return new Response(204);
  });

  server.delete('/plugin-manager/:id', (schema: AppSchema, request) => {
    const { id } = request.params;

    const plugin = schema.find('plugin', id);

    if (!plugin)
      return new Response(404, {}, { code: '404', error: 'Not Found' });

    plugin.destroy();
    return new Response(204, {});
  });
}
