/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BASE_APP: string;
  readonly VITE_BASE_API: string;
  readonly VITE_ENV: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
