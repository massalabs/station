// STYLES

// EXTERNALS
import { persist } from 'zustand/middleware';

// LOCALS

export interface ConfigStoreState {
  lang: string;
  setLang: (lang: string) => void;
  dropLang: () => void;
}

const configStore = persist<ConfigStoreState>(
  (set) => ({
    lang: 'en_US',

    setLang: (lang: string) => set({ lang }),

    dropLang: () => set({ lang: 'en_US' }),
  }),
  {
    name: 'config-store',
  },
);

export default configStore;
