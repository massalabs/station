// STYLES

// EXTERNALS
import { Theme } from '@massalabs/react-ui-kit';
import { persist } from 'zustand/middleware';

// LOCALS

export interface ConfigStoreState {
  lang: string;
  theme: Theme;

  setLang: (lang: string) => void;
  setTheme: (theme: Theme) => void;

  dropLang: () => void;
  dropTheme: () => void;
}

const configStore = persist<ConfigStoreState>(
  (set) => ({
    lang: 'en_US',
    theme: 'theme-dark',

    setLang: (lang: string) => set({ lang }),
    setTheme: (theme: Theme) => set({ theme }),

    dropLang: () => set({ lang: 'en_US' }),
    dropTheme: () => set({ theme: 'theme-dark' }),
  }),
  {
    name: 'config-store',
  },
);

export default configStore;
