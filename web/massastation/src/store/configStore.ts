// STYLES

// EXTERNALS
import { Theme } from '@massalabs/react-ui-kit';
import { persist } from 'zustand/middleware';
import { DEFAULT_THEME } from '../layouts/LayoutStation/LayoutStation';

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
    theme: DEFAULT_THEME,

    setLang: (lang: string) => set({ lang }),
    setTheme: (theme: Theme) => set({ theme }),

    dropLang: () => set({ lang: 'en_US' }),
    dropTheme: () => set({ theme: DEFAULT_THEME }),
  }),
  {
    name: 'config-store',
  },
);

export default configStore;
