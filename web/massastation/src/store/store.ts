// STYLES

// EXTERNALS
import { create } from 'zustand';

// LOCALS
import configStore, { ConfigStoreState } from './configStore';
import accountStore, { AccountStoreState } from './accountStore';

export const useConfigStore = create<ConfigStoreState>((...obj) => ({
  ...configStore(...obj),
}));

export const useAccountStore = create<AccountStoreState>((...obj) => ({
  ...accountStore(...obj),
}));
