import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { DomainModel } from '../../models/DomainModel';

export const domainFactory = Factory.extend<DomainModel>({
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
  favicon() {
    return 'https://massa.net/favicons/apple-touch-icon.png';
  },
});
