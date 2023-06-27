import { Factory } from 'miragejs';
import { faker } from '@faker-js/faker';
import { AccountObject } from '../../models/AccountModel';

export const accountFactory = Factory.extend<AccountObject>({
  nickname() {
    return faker.internet.userName();
  },
  candidateBalance() {
    return faker.number.int().toString();
  },
  balance() {
    return faker.number.int().toString();
  },
  address() {
    return 'AU' + faker.string.alpha({ length: { min: 128, max: 256 } });
  },
  keyPair() {
    return {
      privateKey: faker.string.alpha({ length: { min: 128, max: 256 } }),
      publicKey: faker.string.alpha({ length: { min: 128, max: 256 } }),
      nonce: faker.lorem.word(),
      salt: faker.lorem.word(),
    };
  },
});
