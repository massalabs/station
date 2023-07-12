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
