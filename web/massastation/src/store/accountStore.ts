// EXTERNALS
import { persist } from 'zustand/middleware';

export interface AccountStoreState {
  currentAccount: string | null;
  setCurrentAccount: (nickname: string | null) => void;
}

const accountStore = persist<AccountStoreState>(
  (set) => ({
    currentAccount: null,
    setCurrentAccount: (nickname: string | null) =>
      set({ currentAccount: nickname }),
  }),
  {
    name: 'account-store',
  },
);

export default accountStore;
