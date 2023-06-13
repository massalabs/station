import { createServer, Model, Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { ENV } from '../';
const generatePlugin = () => {
  const name = faker.lorem.word();
  return {
    description: faker.lorem.sentence(),
    home: `/plugin/massalabs/${name.toLowerCase()}/`,
    id: faker.number.int(),
    logo: `/plugin/massalabs/${name.toLowerCase()}/logo.svg`,
    name,
    status: 'Up',
    version: faker.system.semver(),
  };
};

/**
 *
 * @param environment
 */
// eslint-disable-next-line jsdoc/require-jsdoc
function mockServer(environment = ENV.DEV) {
  const mockedServer = createServer({
    environment,
    models: {
      plugin: Model,
    },
    factories: {
      plugin: Factory.extend({
        plugin: generatePlugin(),
      }),
    },
    seeds(server) {
      server.createList('plugin', 5);
    },
    routes() {
      this.namespace = import.meta.env.VITE_BASE_API;
      this.get('/plugin-manager', (schema) => {
        let { models: plugins } = schema.plugins.all();
        return plugins;
      });
    },
  });

  return mockedServer;
}

export default mockServer;
