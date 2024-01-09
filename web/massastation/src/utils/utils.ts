export function routeFor(path: string) {
  return `${import.meta.env.VITE_BASE_APP}/${path}`;
}
