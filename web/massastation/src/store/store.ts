// STYLES

// EXTERNALS
import { create } from 'zustand';

// LOCALS
import configStore, { ConfigStoreState } from './configStore';

export const useConfigStore = create<ConfigStoreState>((...obj) => ({
  ...configStore(...obj),
}));
