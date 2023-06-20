import { createServer, Model, Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { ENV } from '../const/env/env';

/**
 *
 * @param environment
 */
function mockServer(environment = ENV.DEV) {
  const commonVariable = 'Common Value';

  const mockedServer = createServer({
    environment,
    models: {
      plugin: Model,
    },
    factories: {
      plugin: Factory.extend({
        name() {
          return faker.lorem.word();
        },
        author() {
          // there is a 30% chance that the author will be MassaLabs
          return Math.random() < 0.3 ? 'Massa Labs' : faker.person.firstName();
        },
        description() {
          return faker.lorem.sentence();
        },
        home() {
          const name = this.name.toLowerCase();
          return `/plugin/massalabs/${name}/`;
        },
        logo() {
          const name = this.name.toLowerCase();
          return `/plugin/massalabs/${name}/logo.svg`;
        },
        status: 'Up',
        version() {
          return faker.system.semver();
        },
        updatable() {
          return Math.random() < 0.5;
        },
      }),
    },
    seeds(server) {
      server.createList('plugin', 2);
    },
    routes() {
      this.get('/plugin-manager', (schema) => {
        let { models: plugins } = schema.plugins.all();
        return plugins;
      });
      this.get('/plugin-manager/:id', (schema, request) => {
        const { id } = request.params;
        let plugin = schema.plugins.find(id);

        if (!plugin)
          return new Response(404, {}, { code: '404', error: 'Not Found' });

        return plugin.attrs;
      });

      this.post('/plugin-manager/:id/execute', (schema, request) => {
        const { command } = JSON.parse(request.requestBody);
        const { id } = request.params;
        const plugin = schema.plugins.find(id);

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

      this.delete('/plugin-manager/:id', (schema, request) => {
        const plugin = schema.plugins.find(request.params.id);

        if (!plugin)
          return new Response(404, {}, { code: '404', error: 'Not Found' });

        return plugin.destroy();
      });
    },
  });

  return mockedServer;
}

export default mockServer;
