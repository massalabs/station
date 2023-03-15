import { Store } from "pullstate";

interface IUIStore {
  theme: string;
}

export const UIStore = new Store<IUIStore>({
  theme: "light",
});