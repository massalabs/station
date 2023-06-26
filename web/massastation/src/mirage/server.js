import { createServer, Model, Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { ENV } from '../const/env/env';

/**
 * Creates a mocked server
 *
 * @param environment the environment to mock
 * @returns the mocked server
 */
function mockServer(environment = ENV.DEV) {
  const mockedServer = createServer({
    environment,
    models: {
      plugin: Model,
      store: Model,
      domain: Model,
      account: Model,
      website: Model,
      network: Model,
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
      domain: Factory.extend({
        name() {
          return faker.internet.domainName();
        },
        address() {
          return 'AU' + faker.string.alpha({ length: { min: 128, max: 256 } });
        },
        description() {
          return faker.lorem.sentence();
        },
        metadata() {
          return 'test';
        },
      }),
      store: Factory.extend({
        name() {
          return faker.lorem.word();
        },
        author() {
          return Math.random() < 0.3 ? 'Massa Labs' : faker.person.firstName();
        },
        description() {
          return faker.lorem.sentence();
        },
        version() {
          return faker.system.semver();
        },
        url() {
          return faker.internet.url();
        },
        logo: 'logo.png',
        massastationMinVersion() {
          return faker.system.semver();
        },
        file() {
          return {
            url: faker.internet.url(),
            checksum: faker.lorem.word(),
          };
        },
        os: 'linux',
      }),
      account: Factory.extend({
        nickname() {
          return faker.internet.userName();
        },
        candidateBalance() {
          return faker.number.int().toString();
        },
        address() {
          return 'AU' + faker.string.alpha({ length: { min: 128, max: 256 } });
        },
      }),
      website: Factory.extend({
        description() {
          return faker.lorem.sentence();
        },
        name() {
          return faker.lorem.word();
        },
        address() {
          return 'AU' + faker.string.alpha({ length: { min: 128, max: 256 } });
        },
        brokenChunks: [],
      }),
      network: Factory.extend({
        availableNetworks() {
          return ['testnet', 'buildnet', 'labnet'];
        },
        currentNetwork() {
          return 'buildnet';
        },
      }),
    },

    seeds(server) {
      server.createList('plugin', 2);
      server.createList('domain', 50);
      server.createList('store', 7);
      server.createList('account', 5);
      server.createList('website', 2);
      server.create('network');
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

      this.post('plugin-manager', (schema, request) => {
        const sourceURL = request.queryParams.source;

        const storePlugin = schema.stores.findBy(
          (store) => store.file.url === sourceURL.split('?').pop(),
        );
        const newPlugin = {
          id: faker.number.int(),
          status: 'Up',
          updatable: false,
          name: storePlugin?.name || faker.lorem.word(),
          author: storePlugin?.author || faker.person.firstName(),
          description: storePlugin?.description || faker.lorem.sentence(),
          version: storePlugin?.version || faker.system.semver(),
          home: storePlugin?.url || `/plugin/massalabs/${faker.lorem.word()}/`,
        };
        schema.plugins.create(newPlugin);

        return new Response(204);
      });

      this.delete('/plugin-manager/:id', (schema, request) => {
        const plugin = schema.plugins.find(request.params.id);

        if (!plugin)
          return new Response(404, {}, { code: '404', error: 'Not Found' });

        return plugin.destroy();
      });

      this.get('/network', (schema) => {
        let { models: network } = schema.networks.all();

        return network.pop().attrs;
      });

      this.get(
        '/all/domains',
        (schema) => {
          let { models: domains } = schema.domains.all();

          return domains;
        },
        { timing: 5000 },
      );

      this.get('plugin-store', (schema) => {
        let { models: stores } = schema.stores.all();

        return stores;
      });

      this.get(
        'plugin/massalabs/wallet/api/accounts',
        (schema) => {
          let { models: accounts } = schema.accounts.all();

          return accounts;
        },
        { timing: 500 },
      );

      this.put('websiteUploader/prepare', (schema) => {
        // TODO: fix this, it doesn't work as expected, it returns only {id: 4}
        return schema.create('website');
      });
    },
  });

  return mockedServer;
}

export default mockServer;
