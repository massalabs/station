// EXTERNALS
import { persist } from 'zustand/middleware';

export interface NetworktoreState {
  currentNetwork: string | null;
  availableNetworks: string[];
  setCurrentNetwork: (nickname: string | null) => void;
  setAvailableNetworks: (networks: string[]) => void;
}

const networkStore = persist<NetworktoreState>(
  (set) => ({
    currentNetwork: null,
    availableNetworks: [],
    setCurrentNetwork: (nickname: string | null) =>
      set({ currentNetwork: nickname }),
    setAvailableNetworks: (networks: string[]) =>
      set({ availableNetworks: networks }),
  }),
  {
    name: 'network-store',
  },
);

export default networkStore;
