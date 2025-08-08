export interface NetworkModel {
  currentNetwork: string;
  availableNetworkInfos: Array<{
    name: string;
    url: string;
    chainId: number;
    version: string;
    status?: 'up' | 'down';
  }>;
}
