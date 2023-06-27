import { Model } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';
import { AccountObject } from '../../models/AccountModel';
import { DomainModel } from '../../models/DomainModel';
import { NetworkModel } from '../../models/NetworkModel';
import { IMassaPlugin } from '../../../../shared/interfaces/IPlugin';
import { IMassaStore } from '../../../../shared/interfaces/IPlugin';

const accountModel: ModelDefinition<AccountObject> = Model.extend({});
const domainModel: ModelDefinition<DomainModel> = Model.extend({});
const networkModel: ModelDefinition<NetworkModel> = Model.extend({});
const pluginModel: ModelDefinition<IMassaPlugin> = Model.extend({});
const storeModel: ModelDefinition<IMassaStore> = Model.extend({});

export const models = {
  account: accountModel,
  domain: domainModel,
  network: networkModel,
  plugin: pluginModel,
  store: storeModel,
};
