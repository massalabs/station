import { createServer, Model, Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { ENV } from '../const/env/env';
const badRequest = new Response(400, {}, { error: 'Bad Request' });
const notFound = new Response(404, {}, { error: 'Not Found' });
const unprocessableEntity = new Response(
  422,
  {},
  { error: 'Unprocessable Entity' },
);

/**
 * Updates the plugin
 *
 * @param plugin the plugin to update
 * @returns the updated plugin
 */
function update(plugin) {
  // increment the version number
  if (plugin === undefined) return unprocessableEntity;
  const pluginVersion = plugin.version.split('.');
  // split the version number into an array
  // increment the patch number
  pluginVersion[2] = parseInt(pluginVersion[2]) + 1;
  plugin.update({
    version: pluginVersion.join('.'),
    updatable: false,
  });
  return plugin;
}
/**
 * Starts the plugin
 *
 * @param plugin the plugin to start
 * @returns the started plugin
 */
function start(plugin) {
  if (plugin === undefined) return unprocessableEntity;
  plugin.update({ status: 'Up' });
  return plugin;
}
/**
 * Stops the plugin
 *
 * @param plugin the plugin to stop
 * @returns the stopped plugin
 */
function stop(plugin) {
  if (plugin === undefined) return unprocessableEntity;
  plugin.update({ status: 'Down' });
  return plugin;
}

/**
 * Creates a mocked server
 *
 * @param environment the environment to mock
 * @returns the mocked server
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
          return Math.random() < 0.3 ? 'MassaLabs' : faker.person.firstName();
        },
        description() {
          return faker.lorem.sentence();
        },
        home() {
          const name = this.name.toLowerCase();
          return `/plugin/massalabs/${name}/`;
        },
        id() {
          return faker.number.int();
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
      server.createList('plugin', Math.floor(Math.random() * 8));
    },
    routes() {
      this.get('/plugin-manager', (schema) => {
        let { models: plugins } = schema.plugins.all();
        return plugins;
      });

      this.post('/plugin-manager/:id/execute', (schema, request) => {
        const { method } = JSON.parse(request.requestBody);
        const { id } = request.params;

        const plugin = schema.plugins.find(id);

        if (!plugin)
          return new Response(404, {}, { code: '404', error: 'Not Found' });

        const status = ['update', 'start'].includes(method) ? 'Up' : 'Down';
        const updatable = method === 'update' ? false : plugin.updatable;

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
