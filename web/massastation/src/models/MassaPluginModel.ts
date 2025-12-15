export interface MassaPluginModel {
  id: string;
  name: string;
  author: string;
  description: string;
  logo: string;
  status: string;
  home: string;
  updatable?: boolean;
  version?: string;
}

// Logical types of plugins managed by MassaLabs.
// String values are aligned with the human‑readable plugin names.
export enum MassaLabsPlugins {
  MassaWallet = 'Massa Wallet',
  NodeManager = 'Node Manager',
  Deweb = 'Deweb',
  // Add additional plugin types here as they are introduced.
}
