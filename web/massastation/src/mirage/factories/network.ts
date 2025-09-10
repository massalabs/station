import { Factory } from 'miragejs';
import { NetworkModel } from '../../models/NetworkModel';

export const networkFactory = Factory.extend<NetworkModel>({
  availableNetworkInfos() {
    return [
      { 
        name: 'mainnet', 
        url: 'https://mainnet.massa.net/api/v2', 
        chainId: 77658377, 
        version: '1.0.0', 
        status: 'up' 
      },
      { 
        name: 'buildnet', 
        url: 'https://buildnet.massa.net/api/v2', 
        chainId: 77658366, 
        version: '1.0.0', 
        status: 'up' 
      },
      { 
        name: 'testnet', 
        url: 'https://testnet.massa.net/api/v2', 
        chainId: 77658355, 
        version: '0.9.0', 
        status: 'down' 
      },
      { 
        name: 'localnet', 
        url: 'http://localhost:33035', 
        chainId: 77658300, 
        version: '1.0.0-dev', 
        status: 'up' 
      },
      { 
        name: 'custom-dev', 
        url: 'https://dev.custom-massa.io/api/v2', 
        chainId: 77658400, 
        version: '1.1.0-beta', 
        status: 'up' 
      }
    ];
  },
  currentNetwork() {
    return 'buildnet';
  },
});
