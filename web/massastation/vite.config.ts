import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import svgr from 'vite-plugin-svgr';
import * as path from 'path';

export default ({ mode }) => {
  // loadEnv(mode, process.cwd()) will load the .env files depending on the mode
  // import.meta.env.VITE_BASE_APP available here with: process.env.VITE_BASE_APP
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };

  return defineConfig({
    plugins: [react(), svgr()],
    base: process.env.VITE_BASE_APP,
    build: {
      emptyOutDir: true,
      manifest: true,
      sourcemap: true,
    },
    resolve: {
      alias: [
        { find: '@', replacement: path.resolve(__dirname, 'src') },
        { find: '@/shared', replacement: path.resolve(__dirname, '../shared') },
        { find: '@wailsjs', replacement: path.resolve(__dirname, 'wailsjs') },
      ],
    },
  });
};
