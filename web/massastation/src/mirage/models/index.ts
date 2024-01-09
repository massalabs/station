import { Model } from 'miragejs';
import { ModelDefinition } from 'miragejs/-types';
import { AccountObject } from '../../models/AccountModel';
import { NetworkModel } from '../../models/NetworkModel';
import { MassaPluginModel, MassaStoreModel } from '@/models';

const accountModel: ModelDefinition<AccountObject> = Model.extend({});
const networkModel: ModelDefinition<NetworkModel> = Model.extend({});
const pluginModel: ModelDefinition<MassaPluginModel> = Model.extend({});
const storeModel: ModelDefinition<MassaStoreModel> = Model.extend({});

export const models = {
  account: accountModel,
  network: networkModel,
  plugin: pluginModel,
  store: storeModel,
};
