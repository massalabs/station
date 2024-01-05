import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { MassaStoreModel } from '@/models';

export const storeFactory = Factory.extend<MassaStoreModel>({
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
  file() {
    return {
      url: faker.internet.url(),
      checksum: faker.lorem.word(),
    };
  },
  os() {
    return 'Linux';
  },
  isCompatible() {
    return Math.random() < 0.3 ? false : true;
  },
  massastationMinVersion() {
    return '<=' + faker.system.semver();
  },
});
