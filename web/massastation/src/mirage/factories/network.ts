import { Factory } from 'miragejs';
import { NetworkModel } from '../../models/NetworkModel';

export const networkFactory = Factory.extend<NetworkModel>({
  availableNetworks() {
    return ['mainnet', 'buildnet'];
  },
  currentNetwork() {
    return 'buildnet';
  },
});
