// EXTERNALS
import { persist } from 'zustand/middleware';
import type { NetworkInfo } from '@/models/NetworkModel';

export interface NetworkStoreState {
  currentNetwork: string | null;
  setCurrentNetwork: (nickname: string | null) => void;
  availableNetworks: NetworkInfo[];
  setAvailableNetworks: (networks: NetworkInfo[]) => void;
  getChainId: () => bigint | undefined;
}

const networkStore = persist<NetworkStoreState>(
  (set, get) => ({
    currentNetwork: null,
    setCurrentNetwork: (nickname: string | null) =>
      set({ currentNetwork: nickname }),
    availableNetworks: [],
    setAvailableNetworks: (networks: NetworkInfo[]) => set({ availableNetworks: networks }),
    getChainId: (): bigint | undefined => {
      const tempChId = get().availableNetworks.find((n: NetworkInfo) => n.name === get().currentNetwork)?.chainId;
      return tempChId ? BigInt(tempChId) : undefined;
    },
  }),
  {
    name: 'network-store',
  },
);

export default networkStore;
