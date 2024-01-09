// STYLES

// EXTERNALS
import { create } from 'zustand';

// LOCALS
import configStore, { ConfigStoreState } from './configStore';
import networkStore, { NetworkStoreState } from './networkStore';

export const useConfigStore = create<ConfigStoreState>((...obj) => ({
  ...configStore(...obj),
}));

export const useNetworkStore = create<NetworkStoreState>((...obj) => ({
  ...networkStore(...obj),
}));
