// EXTERNALS
import { persist } from 'zustand/middleware';

export interface NetworkStoreState {
  currentNetwork: string | null;
  setCurrentNetwork: (nickname: string | null) => void;
}

const networkStore = persist<NetworkStoreState>(
  (set) => ({
    currentNetwork: null,
    setCurrentNetwork: (nickname: string | null) =>
      set({ currentNetwork: nickname }),
  }),
  {
    name: 'network-store',
  },
);

export default networkStore;
