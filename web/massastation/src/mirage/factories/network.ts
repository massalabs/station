import { Factory } from 'miragejs';
import { NetworkModel } from '../../models/NetworkModel';

export const networkFactory = Factory.extend<NetworkModel>({
  availableNetworkInfos() {
    return [
      { name: 'mainnet', url: 'https://mainnet.massa.net/api/v2', chainId: 77658377, version: '1.0.0', status: 'up' },
      { name: 'buildnet', url: 'https://buildnet.massa.net/api/v2', chainId: 77658366, version: '1.0.0', status: 'up' },
    ];
  },
  currentNetwork() {
    return 'buildnet';
  },
});
