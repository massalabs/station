// STYLES

// EXTERNALS
import { create } from 'zustand';

// LOCALS
import configStore, { ConfigStoreState } from './configStore';
import accountStore, { AccountStoreState } from './accountStore';
import networkStore, { NetworktoreState } from './networkStore';

export const useConfigStore = create<ConfigStoreState>((...obj) => ({
  ...configStore(...obj),
}));

export const useAccountStore = create<AccountStoreState>((...obj) => ({
  ...accountStore(...obj),
}));

export const useNetworkStore = create<NetworktoreState>((...obj) => ({
  ...networkStore(...obj),
}));
