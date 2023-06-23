// EXTERNALS
import { persist } from 'zustand/middleware';

export interface AccountStoreState {
  nickname: string | null;
  setNickname: (nickname: string | null) => void;
}

const accountStore = persist<AccountStoreState>(
  (set) => ({
    nickname: null,
    setNickname: (nickname: string | null) => set({ nickname }),
  }),
  {
    name: 'account-store',
  },
);

export default accountStore;
