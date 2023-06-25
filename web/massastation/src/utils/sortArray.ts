interface PluginName {
  name: string;
}

export function sortPlugins<T extends PluginName>(array: T[]): T[] {
  return array?.sort((a, b) => a.name.localeCompare(b.name));
}
