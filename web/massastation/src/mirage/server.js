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
      }),
    },
    seeds(server) {
      server.createList('plugin', 5);
    },
    routes() {
      this.get('/plugin-manager', (schema) => {
        let { models: plugins } = schema.plugins.all();
        return plugins;
      });
    },
  });

  return mockedServer;
}

export default mockServer;
