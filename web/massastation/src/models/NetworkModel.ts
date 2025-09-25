export interface NetworkInfo {
  name: string;
  url: string;
  chainId: number;
  version: string;
  status?: 'up' | 'down';
}

export interface NetworkModel {
  currentNetwork: string;
  availableNetworkInfos: Array<NetworkInfo>;
}
